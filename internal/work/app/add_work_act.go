package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type AddWorkActUseCase struct {
	repo  repository.WorkRepository
	audit events.AuditLogger
}

func NewAddWorkActUseCase(r repository.WorkRepository, audit events.AuditLogger) *AddWorkActUseCase {
	return &AddWorkActUseCase{repo: r, audit: audit}
}

func (uc *AddWorkActUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req AddWorkActRequest) (*WorkDetailDTO, error) {
	parsedWorkID, err := uuid.Parse(workID)
	if err != nil {
		return nil, errors.New("ID de trabajo inválido")
	}

	actID, err := uuid.Parse(req.ActID)
	if err != nil {
		return nil, errors.New("ID de acto inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedWorkID)
	if err != nil {
		return nil, err
	}
	if work == nil {
		return nil, errors.New("trabajo no encontrado")
	}

	if !canModifyWork(work, reqCtx) {
		return nil, errors.New("no tienes permisos para modificar este trabajo")
	}

	if work.Status == entities.StatusApproved {
		return nil, errors.New("no se puede modificar un trabajo aprobado")
	}

	if err := uc.repo.AddActs(ctx, work.ID, []uuid.UUID{actID}); err != nil {
		return nil, err
	}

	// Build response
	acts, _ := uc.repo.GetActsByWorkID(ctx, work.ID)
	collabs, _ := uc.repo.GetCollaborators(ctx, work.ID)
	clientName, _ := uc.repo.GetClientNameByID(ctx, work.ClientID)
	clientInfo, _ := uc.repo.GetClientByID(ctx, work.ClientID)
	branchName, _ := uc.repo.GetBranchNameByID(ctx, work.BranchID)

	var drafterName string
	if work.MainDrafterID != nil {
		drafterName, _ = uc.repo.GetUserFullNameByID(ctx, *work.MainDrafterID)
	}

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

		actName, _ := uc.repo.GetActNameByID(ctx, actID)

		details := map[string]interface{}{}
		if work.Folio != nil {
			details["folio"] = *work.Folio
		}
		details["act_name"] = actName

		_ = uc.audit.LogAction(ctx, "ADD_ACT", "WORK", work.ID, reqUUID, details)
	}

	return buildWorkDetail(work, acts, collabs, clientName, branchName, drafterName, clientInfo, dedupReqs, workReqs), nil
}
