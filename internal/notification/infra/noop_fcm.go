package infra

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/ports"
)

// noopFCMSender es un sender que no hace nada.
// Se usa cuando Firebase no esta configurado para que el sistema siga
// creando notificaciones in-app y SSE sin intentar enviar push.
type noopFCMSender struct{}

func (n *noopFCMSender) SendToToken(_ context.Context, _ string, _ ports.PushNotificationPayload) error {
	return nil
}

func (n *noopFCMSender) SendToMultipleTokens(_ context.Context, _ []string, _ ports.PushNotificationPayload) error {
	return nil
}
