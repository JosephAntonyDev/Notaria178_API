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

type WorkRepository interface {
	GetCollaborators(ctx context.Context, workID uuid.UUID) ([]WorkCollaboratorInfo, error)
	GetUsersToNotifyForWork(ctx context.Context, workID uuid.UUID) ([]WorkCollaboratorInfo, error)
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
	fmt.Printf("[DEBUG-NOTIF] Iniciando NotifyNewComment para WorkID: %s, Author: %s\n", input.WorkID, input.CommentAuthor)

	// 1. Obtener todos los usuarios que deben ser notificados (colaboradores + superadmins + proyectista principal)
	usersToNotify, err := uc.workRepo.GetUsersToNotifyForWork(ctx, input.WorkID)
	if err != nil {
		fmt.Printf("[DEBUG-NOTIF] Error GetUsersToNotifyForWork: %v\n", err)
		return fmt.Errorf("error obteniendo usuarios a notificar: %w", err)
	}
	fmt.Printf("[DEBUG-NOTIF] usersToNotify count: %d\n", len(usersToNotify))

	// 2. Filtrar al autor del comentario (no se notifica a sí mismo)
	var recipientUserIDs []uuid.UUID
	for _, user := range usersToNotify {
		if user.UserID != input.CommentAuthor {
			recipientUserIDs = append(recipientUserIDs, user.UserID)
		}
	}

	fmt.Printf("[DEBUG-NOTIF] recipientUserIDs count: %d\n", len(recipientUserIDs))

	// Si no hay destinatarios, no hacer nada
	if len(recipientUserIDs) == 0 {
		fmt.Printf("[DEBUG-NOTIF] Abortando: no hay destinatarios\n")
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

	fmt.Printf("[DEBUG-NOTIF] Creando batch de %d notificaciones\n", len(notifications))
	if err := uc.notifRepo.CreateBatch(ctx, notifications); err != nil {
		fmt.Printf("[DEBUG-NOTIF] Error CreateBatch: %v\n", err)
		// Fallback: Si CreateBatch falla, intentar crear uno por uno
		fmt.Printf("[DEBUG-NOTIF] Intentando fallback Create individual...\n")
		for _, n := range notifications {
			_ = uc.notifRepo.Create(ctx, n)
		}
	} else {
		fmt.Printf("[DEBUG-NOTIF] CreateBatch exitoso\n")
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
