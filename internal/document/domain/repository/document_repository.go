package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/entities"
	"github.com/google/uuid"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *entities.Document) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Document, error)
	GetByWorkID(ctx context.Context, workID uuid.UUID) ([]*entities.Document, error)
}
