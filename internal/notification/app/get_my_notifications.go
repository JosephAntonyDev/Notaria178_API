package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

type GetMyNotificationsUseCase struct {
	repo repository.NotificationRepository
}

func NewGetMyNotificationsUseCase(r repository.NotificationRepository) *GetMyNotificationsUseCase {
	return &GetMyNotificationsUseCase{repo: r}
}

type GetMyNotificationsResult struct {
	Notifications []NotificationDTO `json:"notifications"`
	UnreadCount   int               `json:"unread_count"`
}

func (uc *GetMyNotificationsUseCase) Execute(ctx context.Context, userID string, filters repository.NotificationFilters) (*GetMyNotificationsResult, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("error interno al identificar al usuario")
	}

	notifs, err := uc.repo.ListByUser(ctx, parsedID, filters)
	if err != nil {
		return nil, err
	}

	unread, err := uc.repo.CountUnread(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	notifsDTO := make([]NotificationDTO, 0, len(notifs))
	for _, n := range notifs {
		notifsDTO = append(notifsDTO, ToNotificationDTO(n))
	}

	return &GetMyNotificationsResult{
		Notifications: notifsDTO,
		UnreadCount:   unread,
	}, nil
}
