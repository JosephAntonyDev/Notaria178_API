package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

type MarkAllReadUseCase struct {
	repo repository.NotificationRepository
}

func NewMarkAllReadUseCase(r repository.NotificationRepository) *MarkAllReadUseCase {
	return &MarkAllReadUseCase{repo: r}
}

func (uc *MarkAllReadUseCase) Execute(ctx context.Context, userID string) error {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("error interno al identificar al usuario")
	}

	return uc.repo.MarkAllAsRead(ctx, parsedID)
}
