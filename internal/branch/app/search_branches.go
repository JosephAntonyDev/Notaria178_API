package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/repository"
)

type SearchBranchesUseCase struct {
	repo repository.BranchRepository
}

func NewSearchBranchesUseCase(r repository.BranchRepository) *SearchBranchesUseCase {
	return &SearchBranchesUseCase{repo: r}
}

func (uc *SearchBranchesUseCase) Execute(ctx context.Context, filters repository.BranchFilters) ([]BranchDTO, error) {
	branches, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	branchesDTO := make([]BranchDTO, 0, len(branches))
	for _, b := range branches {
		branchesDTO = append(branchesDTO, ToBranchDTO(b))
	}

	return branchesDTO, nil
}
