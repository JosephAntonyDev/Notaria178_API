package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/google/uuid"
)

type UpdateActRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateActUseCase struct {
	repo  repository.ActRepository
	cache cache.CachePort
}

func NewUpdateActUseCase(r repository.ActRepository, c cache.CachePort) *UpdateActUseCase {
	return &UpdateActUseCase{repo: r, cache: c}
}

func (uc *UpdateActUseCase) Execute(ctx context.Context, actID string, req UpdateActRequest) (*ActDTO, error) {
	parsedID, err := uuid.Parse(actID)
	if err != nil {
		return nil, errors.New("ID de acto inválido")
	}

	act, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	if act == nil {
		return nil, errors.New("acto no encontrado")
	}

	if req.Name != nil && *req.Name != act.Name {
		existing, _ := uc.repo.GetByName(ctx, *req.Name)
		if existing != nil {
			return nil, errors.New("el nombre del acto ya está registrado")
		}
		act.Name = *req.Name
	}

	if req.Description != nil {
		act.Description = req.Description
	}

	if err := uc.repo.Update(ctx, act); err != nil {
		return nil, err
	}

	// Invalidar caché de búsquedas tras mutación
	if uc.cache != nil {
		_ = uc.cache.InvalidatePrefix(ctx, "acts:search:")
	}

	dto := ToActDTO(act)
	return &dto, nil
}
