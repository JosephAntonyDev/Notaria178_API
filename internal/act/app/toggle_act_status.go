package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/google/uuid"
)

type ToggleActStatusUseCase struct {
	repo  repository.ActRepository
	cache cache.CachePort
}

func NewToggleActStatusUseCase(r repository.ActRepository, c cache.CachePort) *ToggleActStatusUseCase {
	return &ToggleActStatusUseCase{repo: r, cache: c}
}

func (uc *ToggleActStatusUseCase) Execute(ctx context.Context, actID string) (*ActDTO, error) {
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

	var newStatus entities.ActStatus
	if act.Status == entities.StatusActive {
		newStatus = entities.StatusInactive
	} else {
		newStatus = entities.StatusActive
	}

	if err := uc.repo.UpdateStatus(ctx, parsedID, newStatus); err != nil {
		return nil, err
	}

	// Invalidar caché de búsquedas tras mutación
	if uc.cache != nil {
		_ = uc.cache.InvalidatePrefix(ctx, "acts:search:")
	}

	act.Status = newStatus
	dto := ToActDTO(act)
	return &dto, nil
}
