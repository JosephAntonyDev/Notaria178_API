package adapters

import (
	"context"

	notifApp "github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	notifEntities "github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	userEntities "github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
	userRepo "github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/google/uuid"
)

// NotifierAdapter adapta *notifApp.CreateNotificationUseCase
// para que cumpla la interfaz work/domain/events.Notifier.
type NotifierAdapter struct {
	uc       *notifApp.CreateNotificationUseCase
	userRepo userRepo.UserRepository
}

func NewNotifierAdapter(uc *notifApp.CreateNotificationUseCase, userRepo userRepo.UserRepository) *NotifierAdapter {
	return &NotifierAdapter{uc: uc, userRepo: userRepo}
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

func (a *NotifierAdapter) NotifySuperAdmins(ctx context.Context, workID *uuid.UUID, notifType string, message string) error {
	if a.userRepo == nil {
		return nil
	}

	role := string(userEntities.RoleSuperAdmin)
	filters := userEntities.UserFilters{
		Role:  &role,
		Limit: 100, // Razonable para el número de notarios/super admins
	}

	admins, _, err := a.userRepo.List(ctx, filters)
	if err != nil {
		return err
	}

	for _, admin := range admins {
		input := notifApp.CreateNotificationInput{
			UserID:  admin.ID,
			WorkID:  workID,
			Type:    notifEntities.NotificationType(notifType),
			Message: message,
		}
		_ = a.uc.Execute(ctx, input)
	}

	return nil
}

// CommentNotifierAdapter adapta *notifApp.NotifyNewCommentUseCase
// para que cumpla la interfaz work/domain/events.CommentNotifier.
type CommentNotifierAdapter struct {
	uc *notifApp.NotifyNewCommentUseCase
}

func NewCommentNotifierAdapter(uc *notifApp.NotifyNewCommentUseCase) *CommentNotifierAdapter {
	return &CommentNotifierAdapter{uc: uc}
}

func (a *CommentNotifierAdapter) NotifyNewComment(ctx context.Context, input events.CommentNotification) error {
	return a.uc.Execute(ctx, notifApp.NotifyNewCommentInput{
		WorkID:         input.WorkID,
		WorkFolio:      input.WorkFolio,
		CommentID:      input.CommentID,
		CommentAuthor:  input.CommentAuthor,
		CommentMessage: input.CommentMessage,
		AuthorName:     input.AuthorName,
	})
}
