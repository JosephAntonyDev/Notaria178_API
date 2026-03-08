package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/repository"
	"github.com/google/uuid"
)

type GetDocumentUseCase struct {
	repo repository.DocumentRepository
}

func NewGetDocumentUseCase(r repository.DocumentRepository) *GetDocumentUseCase {
	return &GetDocumentUseCase{repo: r}
}

func (uc *GetDocumentUseCase) Execute(ctx context.Context, docID string) (*DocumentDTO, error) {
	parsedID, err := uuid.Parse(docID)
	if err != nil {
		return nil, errors.New("ID de documento inválido")
	}

	doc, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("documento no encontrado")
	}

	dto := ToDocumentDTO(doc)
	return &dto, nil
}
