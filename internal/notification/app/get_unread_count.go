package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

type GetUnreadCountUseCase struct {
	repo repository.NotificationRepository
}

func NewGetUnreadCountUseCase(repo repository.NotificationRepository) *GetUnreadCountUseCase {
	return &GetUnreadCountUseCase{
		repo: repo,
	}
}

func (uc *GetUnreadCountUseCase) Execute(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := uc.repo.CountUnread(ctx, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}
