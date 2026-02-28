package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities" 
)


type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.UserStatus) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	List(ctx context.Context, limit int, offset int) ([]*entities.User, error)
}