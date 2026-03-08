package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/google/uuid"
)

type ActFilters struct {
	Limit  int
	Offset int
	Search *string
	Status *string
}

type ActRepository interface {
	Create(ctx context.Context, act *entities.Act) error
	Update(ctx context.Context, act *entities.Act) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.ActStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Act, error)
	GetByName(ctx context.Context, name string) (*entities.Act, error)
	List(ctx context.Context, filters ActFilters) ([]*entities.Act, error)
	
	AddRequirement(ctx context.Context, req *entities.ActRequirement) error
	DeleteRequirement(ctx context.Context, reqID uuid.UUID) error
	DeactivateRequirement(ctx context.Context, reqID uuid.UUID) error
	GetRequirements(ctx context.Context, actID uuid.UUID) ([]*entities.ActRequirement, error)
}
