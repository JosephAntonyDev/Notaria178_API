package app

import (
	"context"
	"errors"
	"strings"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type GetWorkDetailUseCase struct {
	repo repository.WorkRepository
}

func NewGetWorkDetailUseCase(r repository.WorkRepository) *GetWorkDetailUseCase {
	return &GetWorkDetailUseCase{repo: r}
}

func (uc *GetWorkDetailUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string) (*WorkDetailDTO, error) {
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

	// Verificar acceso
	userUUID, _ := uuid.Parse(reqCtx.UserID)
	isCollab, _ := uc.repo.IsCollaborator(ctx, work.ID, userUUID)

	if !CanAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a este trabajo")
	}

	acts, _ := uc.repo.GetActsByWorkID(ctx, work.ID)
	collabs, _ := uc.repo.GetCollaborators(ctx, work.ID)

	clientName, _ := uc.repo.GetClientNameByID(ctx, work.ClientID)
	clientInfo, _ := uc.repo.GetClientByID(ctx, work.ClientID)
	branchName, _ := uc.repo.GetBranchNameByID(ctx, work.BranchID)

	var drafterName string
	if work.MainDrafterID != nil {
		drafterName, _ = uc.repo.GetUserFullNameByID(ctx, *work.MainDrafterID)
	}

	// Deduplicate requirements across all acts
	actIDs := make([]uuid.UUID, 0, len(acts))
	actMap := make(map[uuid.UUID]string, len(acts))
	for _, a := range acts {
		actIDs = append(actIDs, a.ActID)
		actMap[a.ActID] = a.Name
	}

	var allReqs []entities.ActRequirementInfo
	var dedupReqs []DeduplicatedReqDTO
	if len(actIDs) > 0 {
		allReqs, _ = uc.repo.GetRequirementsByActIDs(ctx, actIDs)
		dedupReqs = DeduplicateRequirements(allReqs, actMap)
	}
	if dedupReqs == nil {
		dedupReqs = []DeduplicatedReqDTO{}
	}

	// Get work extra requirements
	workReqs, _ := uc.repo.GetWorkRequirements(ctx, work.ID)

	// Look up uploaded CLIENT_REQUIREMENT documents and match by requirement_id
	reqDocs, _ := uc.repo.GetRequirementDocumentsByWorkID(ctx, work.ID)
	if len(reqDocs) > 0 {
		// For dedup reqs: check all contributing act_requirement IDs
		for i, req := range dedupReqs {
			for _, r := range allReqs {
				if strings.EqualFold(strings.TrimSpace(r.Name), strings.TrimSpace(req.Name)) {
					if docID, ok := reqDocs[r.ID]; ok {
						docIDStr := docID.String()
						dedupReqs[i].DocumentID = &docIDStr
						break
					}
				}
			}
		}
		// For work reqs: match by work_requirement ID
		for i, wr := range workReqs {
			if wr.DocumentID == nil {
				if docID, ok := reqDocs[wr.ID]; ok {
					workReqs[i].DocumentID = &docID
				}
			}
		}
	}

	return buildWorkDetail(work, acts, collabs, clientName, branchName, drafterName, clientInfo, dedupReqs, workReqs), nil
}
