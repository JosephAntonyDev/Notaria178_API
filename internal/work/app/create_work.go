package app

import (
	"context"
	"errors"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type CreateWorkUseCase struct {
	repo repository.WorkRepository
}

func NewCreateWorkUseCase(r repository.WorkRepository) *CreateWorkUseCase {
	return &CreateWorkUseCase{repo: r}
}

func (uc *CreateWorkUseCase) Execute(ctx context.Context, reqCtx RequestContext, req CreateWorkRequest) (*WorkDetailDTO, error) {
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return nil, errors.New("ID de sucursal inválido")
	}

	if reqCtx.UserRole != "SUPER_ADMIN" {
		if reqCtx.BranchID != req.BranchID {
			return nil, errors.New("no puedes crear trabajos en una sucursal que no es la tuya")
		}
	}

	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		return nil, errors.New("ID de cliente inválido")
	}

	actIDs := make([]uuid.UUID, 0, len(req.ActIDs))
	for _, aid := range req.ActIDs {
		parsed, err := uuid.Parse(aid)
		if err != nil {
			return nil, errors.New("uno de los IDs de acto es inválido")
		}
		actIDs = append(actIDs, parsed)
	}

	var mainDrafterID *uuid.UUID
	if req.MainDrafterID != nil && *req.MainDrafterID != "" {
		parsed, err := uuid.Parse(*req.MainDrafterID)
		if err != nil {
			return nil, errors.New("ID de proyectista principal inválido")
		}
		mainDrafterID = &parsed
	} else if reqCtx.UserRole != "SUPER_ADMIN" && reqCtx.UserRole != "LOCAL_ADMIN" {
		parsed, err := uuid.Parse(reqCtx.UserID)
		if err != nil {
			return nil, errors.New("error interno al identificar al usuario")
		}
		mainDrafterID = &parsed
	}

	var deadline *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		parsed, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			return nil, errors.New("formato de fecha límite inválido, usa YYYY-MM-DD")
		}
		deadline = &parsed
	}

	now := time.Now()
	newWork := &entities.Work{
		ID:            uuid.New(),
		BranchID:      branchID,
		ClientID:      clientID,
		MainDrafterID: mainDrafterID,
		Folio:         req.Folio,
		Status:        entities.StatusInProgress,
		Deadline:      deadline,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := uc.repo.Create(ctx, newWork); err != nil {
		return nil, err
	}

	if err := uc.repo.AddActs(ctx, newWork.ID, actIDs); err != nil {
		return nil, err
	}

	acts, _ := uc.repo.GetActsByWorkID(ctx, newWork.ID)
	collabs, _ := uc.repo.GetCollaborators(ctx, newWork.ID)

	return buildWorkDetail(newWork, acts, collabs), nil
}
