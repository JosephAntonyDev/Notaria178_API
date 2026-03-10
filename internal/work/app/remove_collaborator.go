package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type RemoveCollaboratorUseCase struct {
	repo  repository.WorkRepository
	audit events.AuditLogger
}

func NewRemoveCollaboratorUseCase(r repository.WorkRepository, audit events.AuditLogger) *RemoveCollaboratorUseCase {
	return &RemoveCollaboratorUseCase{repo: r, audit: audit}
}

func (uc *RemoveCollaboratorUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, targetUserID string) error {
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

	userID, err := uuid.Parse(targetUserID)
	if err != nil {
		return errors.New("ID de colaborador inválido")
	}

	if err := uc.repo.RemoveCollaborator(ctx, parsedWorkID, userID); err != nil {
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
		details["removed_user"] = collabName

		_ = uc.audit.LogAction(ctx, "UNASSIGN", "WORK", parsedWorkID, reqUUID, details)
	}

	return nil
}
