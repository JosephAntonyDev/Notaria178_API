package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/repository"
	"github.com/google/uuid"
)

type UpdateBranchUseCase struct {
	repo repository.BranchRepository
}

func NewUpdateBranchUseCase(r repository.BranchRepository) *UpdateBranchUseCase {
	return &UpdateBranchUseCase{repo: r}
}

func (uc *UpdateBranchUseCase) Execute(ctx context.Context, branchID string, req UpdateBranchRequest) (*BranchDTO, error) {
	parsedID, err := uuid.Parse(branchID)
	if err != nil {
		return nil, errors.New("ID de sucursal inválido")
	}

	branch, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	if branch == nil {
		return nil, errors.New("sucursal no encontrada")
	}

	if req.Name != nil && *req.Name != branch.Name {
		existing, _ := uc.repo.GetByName(ctx, *req.Name)
		if existing != nil {
			return nil, errors.New("el nombre de la sucursal ya está registrado")
		}
		branch.Name = *req.Name
	}

	if req.Address != nil {
		branch.Address = req.Address
	}

	if err := uc.repo.Update(ctx, branch); err != nil {
		return nil, err
	}

	dto := ToBranchDTO(branch)
	return &dto, nil
}
