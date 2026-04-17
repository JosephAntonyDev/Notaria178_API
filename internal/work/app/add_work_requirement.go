package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type AddWorkRequirementUseCase struct {
	repo  repository.WorkRepository
	audit events.AuditLogger
}

func NewAddWorkRequirementUseCase(r repository.WorkRepository, audit events.AuditLogger) *AddWorkRequirementUseCase {
	return &AddWorkRequirementUseCase{repo: r, audit: audit}
}

func (uc *AddWorkRequirementUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req AddWorkRequirementRequest) (*WorkRequirementDTO, error) {
	parsedID, err := uuid.Parse(workID)
	if err != nil {
		return nil, errors.New("ID de trabajo inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	if work == nil {
		return nil, errors.New("trabajo no encontrado")
	}

	userUUID, _ := uuid.Parse(reqCtx.UserID)
	isCollab, _ := uc.repo.IsCollaborator(ctx, work.ID, userUUID)

	if !CanAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a este trabajo")
	}

	if work.Status == entities.StatusApproved {
		return nil, errors.New("no se puede modificar un trabajo aprobado")
	}

	wr, err := uc.repo.AddWorkRequirement(ctx, work.ID, req.Name)
	if err != nil {
		return nil, errors.New("error al agregar el requisito")
	}

	if uc.audit != nil {
		var reqUUID *uuid.UUID
		if parsed, err := uuid.Parse(reqCtx.UserID); err == nil {
			reqUUID = &parsed
		}

		details := map[string]interface{}{"requirement_name": req.Name}
		if work.Folio != nil {
			details["folio"] = *work.Folio
		}

		_ = uc.audit.LogAction(ctx, "ADD_REQUIREMENT", "WORK", work.ID, reqUUID, details)
	}

	dto := &WorkRequirementDTO{
		ID:         wr.ID,
		WorkID:     wr.WorkID,
		Name:       wr.Name,
		DocumentID: wr.DocumentID,
		CreatedAt:  wr.CreatedAt,
	}
	return dto, nil
}
