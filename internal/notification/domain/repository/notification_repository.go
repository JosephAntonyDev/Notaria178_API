package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/google/uuid"
)

type NotificationFilters struct {
	Limit  int
	Offset int
	IsRead *bool // nil = todas, true = leídas, false = no leídas
}

type NotificationRepository interface {
	Create(ctx context.Context, notif *entities.Notification) error
	ListByUser(ctx context.Context, userID uuid.UUID, filters NotificationFilters) ([]*entities.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int, error)
}
