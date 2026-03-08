package app

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/storage"
	"github.com/google/uuid"
)

type UploadDocumentUseCase struct {
	repo    repository.DocumentRepository
	storage *storage.LocalFileStorage
}

func NewUploadDocumentUseCase(r repository.DocumentRepository, s *storage.LocalFileStorage) *UploadDocumentUseCase {
	return &UploadDocumentUseCase{repo: r, storage: s}
}

type UploadDocumentInput struct {
	File         *multipart.FileHeader
	BranchID     string
	WorkID       string
	ClientID     *string
	UserID       string
	DocumentName string
	Category     string
}

func (uc *UploadDocumentUseCase) Execute(ctx context.Context, input UploadDocumentInput) (*DocumentDTO, error) {
	// Validar categoría
	category := entities.DocumentCategory(input.Category)
	switch category {
	case entities.CategoryDraftDeed, entities.CategoryFinalDeed,
		entities.CategoryClientRequirement, entities.CategoryOther:
	default:
		return nil, errors.New("categoría de documento inválida")
	}

	// Parsear IDs
	workID, err := uuid.Parse(input.WorkID)
	if err != nil {
		return nil, errors.New("ID de trabajo inválido")
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, errors.New("error interno al identificar al usuario")
	}

	var clientID *uuid.UUID
	if input.ClientID != nil && *input.ClientID != "" {
		parsed, err := uuid.Parse(*input.ClientID)
		if err != nil {
			return nil, errors.New("ID de cliente inválido")
		}
		clientID = &parsed
	}

	// Guardar archivo físico en disco
	filePath, err := uc.storage.SaveFile(input.File, input.BranchID, input.WorkID)
	if err != nil {
		return nil, errors.New("error al guardar el archivo en disco")
	}

	// Crear registro en BD
	doc := &entities.Document{
		ID:           uuid.New(),
		ClientID:     clientID,
		WorkID:       &workID,
		UserID:       &userID,
		DocumentName: input.DocumentName,
		Category:     category,
		Version:      1,
		FilePath:     filePath,
		CreatedAt:    time.Now(),
	}

	if err := uc.repo.Create(ctx, doc); err != nil {
		return nil, errors.New("error al registrar el documento en la base de datos")
	}

	dto := ToDocumentDTO(doc)
	return &dto, nil
}
