package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
)

type SearchActsUseCase struct {
	repo repository.ActRepository
}

func NewSearchActsUseCase(r repository.ActRepository) *SearchActsUseCase {
	return &SearchActsUseCase{repo: r}
}

func (uc *SearchActsUseCase) Execute(ctx context.Context, filters repository.ActFilters) ([]ActDTO, error) {
	acts, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	actsDTO := make([]ActDTO, 0)
	for _, a := range acts {
		actsDTO = append(actsDTO, ToActDTO(a))
	}

	return actsDTO, nil
}
