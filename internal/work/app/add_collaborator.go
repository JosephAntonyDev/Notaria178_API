package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type AddCollaboratorUseCase struct {
	repo repository.WorkRepository
}

func NewAddCollaboratorUseCase(r repository.WorkRepository) *AddCollaboratorUseCase {
	return &AddCollaboratorUseCase{repo: r}
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

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.New("ID de colaborador inválido")
	}

	return uc.repo.AddCollaborator(ctx, parsedWorkID, userID)
}
