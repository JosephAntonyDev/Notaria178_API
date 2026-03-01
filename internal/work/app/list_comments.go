package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type ListCommentsUseCase struct {
	repo repository.WorkRepository
}

func NewListCommentsUseCase(r repository.WorkRepository) *ListCommentsUseCase {
	return &ListCommentsUseCase{repo: r}
}

func (uc *ListCommentsUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string) ([]WorkCommentDTO, error) {
	parsedWorkID, err := uuid.Parse(workID)
	if err != nil {
		return nil, errors.New("ID de trabajo inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedWorkID)
	if err != nil {
		return nil, err
	}
	if work == nil {
		return nil, errors.New("trabajo no encontrado")
	}

	userUUID, _ := uuid.Parse(reqCtx.UserID)
	isCollab, _ := uc.repo.IsCollaborator(ctx, work.ID, userUUID)
	if !canAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a los comentarios de este trabajo")
	}

	comments, err := uc.repo.GetCommentsByWorkID(ctx, parsedWorkID)
	if err != nil {
		return nil, err
	}

	commentsDTO := make([]WorkCommentDTO, 0, len(comments))
	for _, c := range comments {
		commentsDTO = append(commentsDTO, ToWorkCommentDTO(c))
	}

	return commentsDTO, nil
}
