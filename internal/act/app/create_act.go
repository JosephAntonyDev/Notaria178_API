package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/google/uuid"
)

type CreateActRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

type CreateActUseCase struct {
	repo  repository.ActRepository
	cache cache.CachePort
}

func NewCreateActUseCase(r repository.ActRepository, c cache.CachePort) *CreateActUseCase {
	return &CreateActUseCase{repo: r, cache: c}
}

func (uc *CreateActUseCase) Execute(ctx context.Context, req CreateActRequest) (*ActDTO, error) {
	existing, _ := uc.repo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("el nombre del acto ya está registrado")
	}

	newAct := &entities.Act{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Status:      entities.StatusActive,
	}

	if err := uc.repo.Create(ctx, newAct); err != nil {
		return nil, err
	}

	// Invalidar caché de búsquedas tras mutación
	if uc.cache != nil {
		_ = uc.cache.InvalidatePrefix(ctx, "acts:search:")
	}

	dto := ToActDTO(newAct)
	return &dto, nil
}
