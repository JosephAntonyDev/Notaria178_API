package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/google/uuid"
)

type DeleteActUseCase struct {
	repo  repository.ActRepository
	cache cache.CachePort
}

func NewDeleteActUseCase(r repository.ActRepository, c cache.CachePort) *DeleteActUseCase {
	return &DeleteActUseCase{repo: r, cache: c}
}

func (uc *DeleteActUseCase) Execute(ctx context.Context, actID string) error {
	parsedID, err := uuid.Parse(actID)
	if err != nil {
		return errors.New("ID de acto inválido")
	}

	act, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return err
	}
	if act == nil {
		return errors.New("acto no encontrado")
	}

	if err := uc.repo.Delete(ctx, parsedID); err != nil {
		return err
	}

	// Invalidar caché de búsquedas tras eliminación
	if uc.cache != nil {
		_ = uc.cache.InvalidatePrefix(ctx, "acts:search:")
	}

	return nil
}
