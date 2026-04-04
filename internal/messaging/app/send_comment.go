package app

import (
	"context"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/domain/events"
	workApp "github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/google/uuid"
)

type SendCommentUseCase struct {
	addCommentUC *workApp.AddCommentUseCase
	broadcaster  events.Broadcaster
}

func NewSendCommentUseCase(ac *workApp.AddCommentUseCase, b events.Broadcaster) *SendCommentUseCase {
	return &SendCommentUseCase{addCommentUC: ac, broadcaster: b}
}

type SendCommentInput struct {
	WorkID  string
	Message string
}

func (uc *SendCommentUseCase) Execute(ctx context.Context, reqCtx workApp.RequestContext, input SendCommentInput) error {
	// Usa el AddCommentUseCase existente para persistir y validar permisos
	comment, err := uc.addCommentUC.Execute(ctx, reqCtx, input.WorkID, workApp.AddCommentRequest{
		Message: input.Message,
	})
	if err != nil {
		return err
	}

	// Broadcast a todos los conectados en la sala
	workUUID, _ := uuid.Parse(input.WorkID)
	room := entities.NewWorkRoom(workUUID)

	wsMsg := &entities.WSMessage{
		Type:   entities.WSTypeComment,
		RoomID: room.ID,
		Payload: entities.CommentPayload{
			ID:        comment.ID,
			WorkID:    comment.WorkID,
			UserID:    comment.UserID,
			UserName:  comment.UserName,
			Message:   comment.Message,
			CreatedAt: comment.CreatedAt,
		},
		Timestamp: time.Now(),
	}

	uc.broadcaster.BroadcastToRoom(room.ID, wsMsg)

	return nil
}
