package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type RemoveCollaboratorUseCase struct {
	repo repository.WorkRepository
}

func NewRemoveCollaboratorUseCase(r repository.WorkRepository) *RemoveCollaboratorUseCase {
	return &RemoveCollaboratorUseCase{repo: r}
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

	userID, err := uuid.Parse(targetUserID)
	if err != nil {
		return errors.New("ID de colaborador inválido")
	}

	return uc.repo.RemoveCollaborator(ctx, parsedWorkID, userID)
}
