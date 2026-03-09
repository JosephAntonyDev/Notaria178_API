package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/storage"
	"github.com/google/uuid"
)

type DeleteDocumentUseCase struct {
	repo    repository.DocumentRepository
	storage *storage.LocalFileStorage
}

func NewDeleteDocumentUseCase(r repository.DocumentRepository, s *storage.LocalFileStorage) *DeleteDocumentUseCase {
	return &DeleteDocumentUseCase{repo: r, storage: s}
}

func (uc *DeleteDocumentUseCase) Execute(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.New("ID de documento inválido")
	}

	// Obtener el documento para conocer la ruta del archivo
	doc, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("documento no encontrado")
	}

	// Eliminar archivo físico del disco
	if doc.FilePath != "" {
		_ = uc.storage.DeleteFile(doc.FilePath)
	}

	// Eliminar registro de la base de datos
	if err := uc.repo.Delete(ctx, id); err != nil {
		return errors.New("error al eliminar el documento de la base de datos")
	}

	return nil
}
