package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type AddCollaboratorUseCase struct {
	repo  repository.WorkRepository
	audit events.AuditLogger
}

func NewAddCollaboratorUseCase(r repository.WorkRepository, audit events.AuditLogger) *AddCollaboratorUseCase {
	return &AddCollaboratorUseCase{repo: r, audit: audit}
}

func (uc *AddCollaboratorUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req AddCollaboratorRequest) error {
	parsedWorkID, err := uuid.Parse(workID)
	if err != nil {
		return errors.New("ID de trabajo inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedWorkID)
	if err != nil {
		return err
	}
	if work == nil {
		return errors.New("trabajo no encontrado")
	}

	if !canModifyWork(work, reqCtx) {
		return errors.New("no tienes permisos para gestionar colaboradores en este trabajo")
	}

	if work.Status == entities.StatusApproved {
		return errors.New("no se puede modificar un trabajo aprobado")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.New("ID de colaborador inválido")
	}

	if err := uc.repo.AddCollaborator(ctx, parsedWorkID, userID); err != nil {
		return err
	}

	if uc.audit != nil {
		var reqUUID *uuid.UUID
		if parsed, err := uuid.Parse(reqCtx.UserID); err == nil {
			reqUUID = &parsed
		}
		
		collabName, _ := uc.repo.GetUserFullNameByID(ctx, userID)
		
		details := map[string]interface{}{}
		if work.Folio != nil {
			details["folio"] = *work.Folio
		}
		details["assigned_to"] = collabName
		
		_ = uc.audit.LogAction(ctx, "ASSIGN", "WORK", parsedWorkID, reqUUID, details)
	}

	return nil
}
