package app

import (
	"context"
	"errors"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type AddCommentUseCase struct {
	repo  repository.WorkRepository
	audit events.AuditLogger
}

func NewAddCommentUseCase(r repository.WorkRepository, audit events.AuditLogger) *AddCommentUseCase {
	return &AddCommentUseCase{repo: r, audit: audit}
}

func (uc *AddCommentUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req AddCommentRequest) (*WorkCommentDTO, error) {
	parsedWorkID, err := uuid.Parse(workID)
	if err != nil {
		return nil, errors.New("ID de trabajo inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedWorkID)
	if err != nil {
		return nil, err
	}
	if work == nil {
		return nil, errors.New("trabajo no encontrado")
	}

	userUUID, err := uuid.Parse(reqCtx.UserID)
	if err != nil {
		return nil, errors.New("error interno al identificar al usuario")
	}

	isCollab, _ := uc.repo.IsCollaborator(ctx, work.ID, userUUID)
	if !canAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a este trabajo para comentar")
	}

	comment := &entities.WorkComment{
		ID:        uuid.New(),
		WorkID:    parsedWorkID,
		UserID:    userUUID,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.AddComment(ctx, comment); err != nil {
		return nil, err
	}

	if uc.audit != nil {
		details := map[string]interface{}{
			"message": req.Message,
		}
		if work.Folio != nil {
			details["folio"] = *work.Folio
		}
		
		_ = uc.audit.LogAction(ctx, "COMMENT", "WORK", parsedWorkID, &userUUID, details)
	}

	dto := ToWorkCommentDTO(*comment)
	return &dto, nil
}
