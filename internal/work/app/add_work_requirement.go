package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type AddWorkRequirementUseCase struct {
	repo repository.WorkRepository
}

func NewAddWorkRequirementUseCase(r repository.WorkRepository) *AddWorkRequirementUseCase {
	return &AddWorkRequirementUseCase{repo: r}
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

	if !canAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a este trabajo")
	}

	if work.Status == entities.StatusApproved {
		return nil, errors.New("no se puede modificar un trabajo aprobado")
	}

	wr, err := uc.repo.AddWorkRequirement(ctx, work.ID, req.Name)
	if err != nil {
		return nil, errors.New("error al agregar el requisito")
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
