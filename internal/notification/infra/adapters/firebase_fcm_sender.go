package adapters

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/ports"
	"google.golang.org/api/option"
)

// FirebaseNotificationSender implementa la interfaz NotificationSender usando Firebase Cloud Messaging
type FirebaseNotificationSender struct {
	client *messaging.Client
}

// NewFirebaseNotificationSender crea una nueva instancia del adaptador FCM
// credentialsPath: ruta al archivo JSON de credenciales de Firebase (service account)
func NewFirebaseNotificationSender(credentialsPath string) (*FirebaseNotificationSender, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error inicializando Firebase app: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creando cliente de Firebase Messaging: %w", err)
	}

	log.Println("[OK] Firebase Cloud Messaging inicializado correctamente")
	return &FirebaseNotificationSender{
		client: client,
	}, nil
}

// SendToToken envía una notificación push a un token específico
func (fcm *FirebaseNotificationSender) SendToToken(ctx context.Context, fcmToken string, payload ports.PushNotificationPayload) error {
	// Convertir data payload a map[string]string (requerido por FCM)
	dataPayload := make(map[string]string)
	for key, value := range payload.Data {
		dataPayload[key] = fmt.Sprintf("%v", value)
	}

	message := &messaging.Message{
		Token: fcmToken,
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data: dataPayload,
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: payload.Title,
				Body:  payload.Body,
				Icon:  "/logo.png", // Icono personalizado para web
				Badge: "/badge.png",
			},
			FCMOptions: &messaging.WebpushFCMOptions{
				Link: "/", // URL para abrir al hacer clic
			},
		},
	}

	response, err := fcm.client.Send(ctx, message)
	if err != nil {
		log.Printf("[ERR] Error enviando notificacion push a token %s: %v", fcmToken[:10]+"...", err)
		return fmt.Errorf("error enviando notificación push: %w", err)
	}

	log.Printf("[OK] Notificacion push enviada exitosamente. Message ID: %s", response)
	return nil
}

// SendToMultipleTokens envía una notificación push a múltiples tokens (batch)
func (fcm *FirebaseNotificationSender) SendToMultipleTokens(ctx context.Context, fcmTokens []string, payload ports.PushNotificationPayload) error {
	if len(fcmTokens) == 0 {
		return nil
	}

	// Convertir data payload a map[string]string
	dataPayload := make(map[string]string)
	for key, value := range payload.Data {
		dataPayload[key] = fmt.Sprintf("%v", value)
	}

	message := &messaging.MulticastMessage{
		Tokens: fcmTokens,
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data: dataPayload,
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: payload.Title,
				Body:  payload.Body,
				Icon:  "/logo.png",
				Badge: "/badge.png",
			},
			FCMOptions: &messaging.WebpushFCMOptions{
				Link: "/",
			},
		},
	}

	response, err := fcm.client.SendEachForMulticast(ctx, message)
	if err != nil {
		log.Printf("[ERR] Error enviando notificaciones push en lote: %v", err)
		return fmt.Errorf("error enviando notificaciones push en lote: %w", err)
	}

	log.Printf("[OK] Notificaciones push enviadas: %d exitosas, %d fallidas", response.SuccessCount, response.FailureCount)

	// Log de tokens que fallaron (para debugging)
	if response.FailureCount > 0 {
		for idx, resp := range response.Responses {
			if !resp.Success {
				log.Printf("[WARN] Token fallido [%d]: %v", idx, resp.Error)
			}
		}
	}

	return nil
}
