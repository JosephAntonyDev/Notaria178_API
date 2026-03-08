package events

import "github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"

// NotificationNotifier define el puerto de salida para emitir notificaciones en tiempo real.
type NotificationNotifier interface {
	Broadcast(userID string, notification *entities.Notification)
}
