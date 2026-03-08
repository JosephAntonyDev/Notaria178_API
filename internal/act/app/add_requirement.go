package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/google/uuid"
)

type AddRequirementRequest struct {
	ActID uuid.UUID `json:"-"`
	Name  string    `json:"name" binding:"required"`
}

type AddRequirementUseCase struct {
	repo  repository.ActRepository
	cache cache.CachePort
}

func NewAddRequirementUseCase(r repository.ActRepository, c cache.CachePort) *AddRequirementUseCase {
	return &AddRequirementUseCase{repo: r, cache: c}
}

func (uc *AddRequirementUseCase) Execute(ctx context.Context, req AddRequirementRequest) (*ActRequirementDTO, error) {
	newReq := &entities.ActRequirement{
		ActID: req.ActID,
		Name:  req.Name,
	}

	if err := uc.repo.AddRequirement(ctx, newReq); err != nil {
		return nil, err
	}

	if uc.cache != nil {
		_ = uc.cache.InvalidatePrefix(ctx, "acts:search:")
	}

	dto := ToActRequirementDTO(newReq)
	return &dto, nil
}
