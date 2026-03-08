package adapters

import (
	"context"

	notifApp "github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	notifEntities "github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/google/uuid"
)

// NotifierAdapter adapta *notifApp.CreateNotificationUseCase
// para que cumpla la interfaz work/domain/events.Notifier.
type NotifierAdapter struct {
	uc *notifApp.CreateNotificationUseCase
}

func NewNotifierAdapter(uc *notifApp.CreateNotificationUseCase) *NotifierAdapter {
	return &NotifierAdapter{uc: uc}
}

func (a *NotifierAdapter) SendNotification(ctx context.Context, userID uuid.UUID, workID *uuid.UUID, notifType string, message string) error {
	input := notifApp.CreateNotificationInput{
		UserID:  userID,
		WorkID:  workID,
		Type:    notifEntities.NotificationType(notifType),
		Message: message,
	}
	return a.uc.Execute(ctx, input)
}
