package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/entities"
	"github.com/google/uuid"
)

type ClientFilters struct {
	Limit  int
	Offset int
	Search *string
}

type ClientRepository interface {
	Create(ctx context.Context, client *entities.Client) error
	Update(ctx context.Context, client *entities.Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Client, error)
	GetByRFC(ctx context.Context, rfc string) (*entities.Client, error)
	List(ctx context.Context, filters ClientFilters) ([]*entities.Client, error)
}
