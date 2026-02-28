package app

import (
	"context"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

type CreateNotificationUseCase struct {
	repo     repository.NotificationRepository
	notifier events.NotificationNotifier
}

func NewCreateNotificationUseCase(r repository.NotificationRepository, n events.NotificationNotifier) *CreateNotificationUseCase {
	return &CreateNotificationUseCase{repo: r, notifier: n}
}

type CreateNotificationInput struct {
	UserID  uuid.UUID
	WorkID  *uuid.UUID
	Type    entities.NotificationType
	Message string
}

func (uc *CreateNotificationUseCase) Execute(ctx context.Context, input CreateNotificationInput) error {
	notif := &entities.Notification{
		ID:        uuid.New(),
		UserID:    input.UserID,
		WorkID:    input.WorkID,
		Type:      input.Type,
		Message:   input.Message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, notif); err != nil {
		return err
	}

	// Emitir en tiempo real al usuario conectado
	uc.notifier.Broadcast(input.UserID.String(), notif)

	return nil
}
