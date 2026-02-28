package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
)

type GetProfileUseCase struct {
	repo repository.UserRepository
}

func NewGetProfileUseCase(r repository.UserRepository) *GetProfileUseCase {
	return &GetProfileUseCase{
		repo: r,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID string) (*UserPublicDTO, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("ID de usuario inválido en el token")
	}

	user, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("el usuario no existe en la base de datos")
	}

	profileDTO := ToUserPublicDTO(user)
	return &profileDTO, nil
}