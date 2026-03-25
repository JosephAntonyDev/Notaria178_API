package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/ports"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

// WorkRepository interface simplificada (solo lo que necesitamos)
type WorkRepository interface {
	GetCollaborators(ctx context.Context, workID uuid.UUID) ([]WorkCollaboratorInfo, error)
}

type WorkCollaboratorInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
}

type NotifyNewCommentInput struct {
	WorkID         uuid.UUID `json:"work_id"`
	WorkFolio      string    `json:"work_folio"` // Para el mensaje de notificación
	CommentID      uuid.UUID `json:"comment_id"`
	CommentAuthor  uuid.UUID `json:"comment_author"` // El que escribió el comentario
	CommentMessage string    `json:"comment_message"`
	AuthorName     string    `json:"author_name"` // Nombre del autor
}

type NotifyNewCommentUseCase struct {
	notifRepo   repository.NotificationRepository
	deviceRepo  repository.DeviceTokenRepository
	workRepo    WorkRepository
	fcmSender   ports.NotificationSender
	sseNotifier SSENotifier // Para enviar notificaciones SSE en tiempo real
}

// SSENotifier interface para integracion con SSE existente
type SSENotifier interface {
	Broadcast(userID string, notification *entities.Notification)
}

func NewNotifyNewCommentUseCase(
	notifRepo repository.NotificationRepository,
	deviceRepo repository.DeviceTokenRepository,
	workRepo WorkRepository,
	fcmSender ports.NotificationSender,
	sseNotifier SSENotifier,
) *NotifyNewCommentUseCase {
	return &NotifyNewCommentUseCase{
		notifRepo:   notifRepo,
		deviceRepo:  deviceRepo,
		workRepo:    workRepo,
		fcmSender:   fcmSender,
		sseNotifier: sseNotifier,
	}
}

func (uc *NotifyNewCommentUseCase) Execute(ctx context.Context, input NotifyNewCommentInput) error {
	// 1. Obtener todos los colaboradores del trabajo
	collaborators, err := uc.workRepo.GetCollaborators(ctx, input.WorkID)
	if err != nil {
		return fmt.Errorf("error obteniendo colaboradores: %w", err)
	}

	// 2. Filtrar al autor del comentario (no se notifica a sí mismo)
	var recipientUserIDs []uuid.UUID
	for _, collab := range collaborators {
		if collab.UserID != input.CommentAuthor {
			recipientUserIDs = append(recipientUserIDs, collab.UserID)
		}
	}

	// Si no hay destinatarios, no hacer nada
	if len(recipientUserIDs) == 0 {
		return nil
	}

	// 3. Crear título y cuerpo de la notificación
	title := fmt.Sprintf("Nuevo comentario en %s", input.WorkFolio)
	body := fmt.Sprintf("%s comentó: %s", input.AuthorName, truncateString(input.CommentMessage, 100))
	message := fmt.Sprintf("Nuevo comentario de %s en el trabajo %s", input.AuthorName, input.WorkFolio)

	// 4. Crear notificaciones in-app en lote
	var notifications []*entities.Notification
	for _, userID := range recipientUserIDs {
		titleCopy := title
		bodyCopy := body

		notif := &entities.Notification{
			ID:        uuid.New(),
			UserID:    userID,
			WorkID:    &input.WorkID,
			Type:      entities.TypeNewComment,
			Title:     &titleCopy,
			Body:      &bodyCopy,
			Message:   message,
			IsRead:    false,
			CreatedAt: time.Now(),
		}
		notifications = append(notifications, notif)
	}

	if err := uc.notifRepo.CreateBatch(ctx, notifications); err != nil {
		return fmt.Errorf("error creando notificaciones in-app: %w", err)
	}

	// 5. Enviar notificaciones SSE (en tiempo real para usuarios conectados)
	for _, notif := range notifications {
		if uc.sseNotifier != nil {
			uc.sseNotifier.Broadcast(notif.UserID.String(), notif)
		}
	}

	// 6. Obtener tokens FCM de los destinatarios
	deviceTokens, err := uc.deviceRepo.GetTokensByUserIDs(ctx, recipientUserIDs)
	if err != nil {
		return fmt.Errorf("error obteniendo tokens FCM: %w", err)
	}

	// Si no hay tokens, terminar (los usuarios no tienen dispositivos registrados)
	if len(deviceTokens) == 0 {
		return nil
	}

	// 7. Preparar tokens para FCM
	var fcmTokens []string
	for _, dt := range deviceTokens {
		fcmTokens = append(fcmTokens, dt.FCMToken)
	}

	// 8. Preparar data payload con el comentario completo para actualizar el frontend
	commentData := map[string]interface{}{
		"type":           "NEW_COMMENT",
		"work_id":        input.WorkID.String(),
		"comment_id":     input.CommentID.String(),
		"comment_author": input.CommentAuthor.String(),
		"author_name":    input.AuthorName,
		"message":        input.CommentMessage,
		"created_at":     time.Now().Format(time.RFC3339),
	}

	commentJSON, _ := json.Marshal(commentData)

	payload := ports.PushNotificationPayload{
		Title: title,
		Body:  body,
		Data: map[string]interface{}{
			"click_action": "WORK_DETAIL",
			"work_id":      input.WorkID.String(),
			"comment":      string(commentJSON),
		},
	}

	// 9. Enviar notificaciones push vía FCM
	if err := uc.fcmSender.SendToMultipleTokens(ctx, fcmTokens, payload); err != nil {
		// No retornamos error aquí porque las notificaciones in-app ya se crearon
		// Solo logueamos el error
		fmt.Printf("[WARN] Error enviando notificaciones push (in-app creadas exitosamente): %v\n", err)
	}

	return nil
}

// truncateString trunca un string al límite especificado
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
