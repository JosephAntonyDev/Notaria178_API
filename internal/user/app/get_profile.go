package app

import (
	"context"
	"errors"

	branchEntities "github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/entities"
	branchRepository "github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
	"github.com/google/uuid"
)

type GetProfileUseCase struct {
	repo       repository.UserRepository
	branchRepo branchRepository.BranchRepository
}

func NewGetProfileUseCase(r repository.UserRepository, br branchRepository.BranchRepository) *GetProfileUseCase {
	return &GetProfileUseCase{
		repo:       r,
		branchRepo: br,
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
	var branch *branchEntities.Branch
	if user.BranchID != nil {
		b, err := uc.branchRepo.GetByID(ctx, *user.BranchID)
		if err == nil && b != nil {
			branch = b
		}
	}

	profileDTO := ToUserPublicDTO(user, branch)
	return &profileDTO, nil
}