package app

import (
	"context"
	"errors"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/repository"
	"github.com/google/uuid"
)

type CreateBranchUseCase struct {
	repo repository.BranchRepository
}

func NewCreateBranchUseCase(r repository.BranchRepository) *CreateBranchUseCase {
	return &CreateBranchUseCase{repo: r}
}

func (uc *CreateBranchUseCase) Execute(ctx context.Context, req CreateBranchRequest) (*BranchDTO, error) {
	existing, _ := uc.repo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("el nombre de la sucursal ya está registrado")
	}

	newBranch := &entities.Branch{
		ID:        uuid.New(),
		Name:      req.Name,
		Address:   req.Address,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, newBranch); err != nil {
		return nil, err
	}

	dto := ToBranchDTO(newBranch)
	return &dto, nil
}
