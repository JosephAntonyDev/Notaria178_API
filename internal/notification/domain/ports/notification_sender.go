package ports

import (
	"context"
)

// ─── Push Notification Payload ──────────────────────────────────────────────

type PushNotificationPayload struct {
	Title string                 `json:"title"`
	Body  string                 `json:"body"`
	Data  map[string]interface{} `json:"data,omitempty"` // Data payload para el frontend
}

// ─── Port para enviar notificaciones push ───────────────────────────────────

type NotificationSender interface {
	// Enviar notificación push a un token específico
	SendToToken(ctx context.Context, fcmToken string, payload PushNotificationPayload) error

	// Enviar notificación push a múltiples tokens (batch)
	SendToMultipleTokens(ctx context.Context, fcmTokens []string, payload PushNotificationPayload) error
}
