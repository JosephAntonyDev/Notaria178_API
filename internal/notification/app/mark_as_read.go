package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

type MarkAsReadUseCase struct {
	repo repository.NotificationRepository
}

func NewMarkAsReadUseCase(r repository.NotificationRepository) *MarkAsReadUseCase {
	return &MarkAsReadUseCase{repo: r}
}

func (uc *MarkAsReadUseCase) Execute(ctx context.Context, notifID string, userID string) error {
	parsedNotifID, err := uuid.Parse(notifID)
	if err != nil {
		return errors.New("ID de notificación inválido")
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("error interno al identificar al usuario")
	}

	return uc.repo.MarkAsRead(ctx, parsedNotifID, parsedUserID)
}
