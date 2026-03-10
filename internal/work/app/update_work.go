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

type UpdateWorkUseCase struct {
	repo  repository.WorkRepository
	audit events.AuditLogger
}

func NewUpdateWorkUseCase(r repository.WorkRepository, audit events.AuditLogger) *UpdateWorkUseCase {
	return &UpdateWorkUseCase{repo: r, audit: audit}
}

func (uc *UpdateWorkUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req UpdateWorkRequest) (*WorkDetailDTO, error) {
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

	if !canModifyWork(work, reqCtx) {
		return nil, errors.New("no tienes permisos para modificar este trabajo")
	}

	// Regla Maestra: bloquear mutaciones si el trabajo está APPROVED
	if work.Status == entities.StatusApproved {
		return nil, errors.New("no se puede modificar un trabajo aprobado")
	}

	// Actualizar folio
	if req.Folio != nil {
		work.Folio = req.Folio
	}

	// Actualizar fecha límite
	if req.Deadline != nil && *req.Deadline != "" {
		parsed, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			return nil, errors.New("formato de fecha límite inválido, usa YYYY-MM-DD")
		}
		work.Deadline = &parsed
	}

	if err := uc.repo.Update(ctx, work); err != nil {
		return nil, err
	}

	// Reemplazar actos si se proporcionaron
	if len(req.ActIDs) > 0 {
		actIDs := make([]uuid.UUID, 0, len(req.ActIDs))
		for _, aid := range req.ActIDs {
			parsed, err := uuid.Parse(aid)
			if err != nil {
				return nil, errors.New("uno de los IDs de acto es inválido")
			}
			actIDs = append(actIDs, parsed)
		}
		if err := uc.repo.RemoveAllActs(ctx, work.ID); err != nil {
			return nil, err
		}
		if err := uc.repo.AddActs(ctx, work.ID, actIDs); err != nil {
			return nil, err
		}
	}

	// Construir respuesta
	acts, _ := uc.repo.GetActsByWorkID(ctx, work.ID)
	collabs, _ := uc.repo.GetCollaborators(ctx, work.ID)
	clientName, _ := uc.repo.GetClientNameByID(ctx, work.ClientID)
	clientInfo, _ := uc.repo.GetClientByID(ctx, work.ClientID)
	branchName, _ := uc.repo.GetBranchNameByID(ctx, work.BranchID)

	var drafterName string
	if work.MainDrafterID != nil {
		drafterName, _ = uc.repo.GetUserFullNameByID(ctx, *work.MainDrafterID)
	}

	// Deduplicate requirements
	actIDs := make([]uuid.UUID, 0, len(acts))
	actMap := make(map[uuid.UUID]string, len(acts))
	for _, a := range acts {
		actIDs = append(actIDs, a.ActID)
		actMap[a.ActID] = a.Name
	}
	var dedupReqs []DeduplicatedReqDTO
	if len(actIDs) > 0 {
		allReqs, _ := uc.repo.GetRequirementsByActIDs(ctx, actIDs)
		dedupReqs = DeduplicateRequirements(allReqs, actMap)
	}
	if dedupReqs == nil {
		dedupReqs = []DeduplicatedReqDTO{}
	}
	workReqs, _ := uc.repo.GetWorkRequirements(ctx, work.ID)

	if uc.audit != nil {
		var reqUUID *uuid.UUID
		if parsed, err := uuid.Parse(reqCtx.UserID); err == nil {
			reqUUID = &parsed
		}

		details := map[string]interface{}{}
		if req.Folio != nil {
			details["new_folio"] = *req.Folio
		}
		if req.Deadline != nil {
			details["new_deadline"] = *req.Deadline
		}

		_ = uc.audit.LogAction(ctx, "UPDATE", "WORK", work.ID, reqUUID, details)
	}

	return buildWorkDetail(work, acts, collabs, clientName, branchName, drafterName, clientInfo, dedupReqs, workReqs), nil
}
