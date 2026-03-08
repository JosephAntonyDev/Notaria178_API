package app

import (
	"context"
	"errors"

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

	if !canAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a este trabajo")
	}

	acts, _ := uc.repo.GetActsByWorkID(ctx, work.ID)
	collabs, _ := uc.repo.GetCollaborators(ctx, work.ID)

	return buildWorkDetail(work, acts, collabs), nil
}
