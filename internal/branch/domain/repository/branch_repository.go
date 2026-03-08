package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/entities"
	"github.com/google/uuid"
)

type BranchFilters struct {
	Limit  int
	Offset int
	Search *string
}

type BranchRepository interface {
	Create(ctx context.Context, branch *entities.Branch) error
	Update(ctx context.Context, branch *entities.Branch) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Branch, error)
	GetByName(ctx context.Context, name string) (*entities.Branch, error)
	List(ctx context.Context, filters BranchFilters) ([]*entities.Branch, error)
}
