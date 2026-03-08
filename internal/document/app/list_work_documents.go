package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/repository"
	"github.com/google/uuid"
)

type ListWorkDocumentsUseCase struct {
	repo repository.DocumentRepository
}

func NewListWorkDocumentsUseCase(r repository.DocumentRepository) *ListWorkDocumentsUseCase {
	return &ListWorkDocumentsUseCase{repo: r}
}

func (uc *ListWorkDocumentsUseCase) Execute(ctx context.Context, workID string) ([]DocumentDTO, error) {
	parsedID, err := uuid.Parse(workID)
	if err != nil {
		return nil, errors.New("ID de trabajo inválido")
	}

	docs, err := uc.repo.GetByWorkID(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	docsDTO := make([]DocumentDTO, 0, len(docs))
	for _, d := range docs {
		docsDTO = append(docsDTO, ToDocumentDTO(d))
	}

	return docsDTO, nil
}
