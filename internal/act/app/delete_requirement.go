package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/google/uuid"
)

type DeleteRequirementUseCase struct {
	repo  repository.ActRepository
	cache cache.CachePort
}

func NewDeleteRequirementUseCase(r repository.ActRepository, c cache.CachePort) *DeleteRequirementUseCase {
	return &DeleteRequirementUseCase{repo: r, cache: c}
}

// Execute aplica la lógica de "Eliminación Segura":
// Si el acto padre tiene trabajos vinculados (works_count > 0), desactiva el requisito.
// Si no, lo elimina físicamente.
func (uc *DeleteRequirementUseCase) Execute(ctx context.Context, actID uuid.UUID, reqID uuid.UUID) (softDeleted bool, err error) {
	// Obtener el acto padre para verificar si tiene trabajos vinculados
	act, err := uc.repo.GetByID(ctx, actID)
	if err != nil {
		return false, err
	}

	if act != nil && act.WorksCount > 0 {
		// Desactivar (soft delete) — el acto tiene trabajos vinculados
		if err := uc.repo.DeactivateRequirement(ctx, reqID); err != nil {
			return false, err
		}
		softDeleted = true
	} else {
		// Eliminar físicamente — el acto no tiene trabajos vinculados
		if err := uc.repo.DeleteRequirement(ctx, reqID); err != nil {
			return false, err
		}
		softDeleted = false
	}

	// Invalidar caché
	if uc.cache != nil {
		_ = uc.cache.InvalidatePrefix(ctx, "acts:search:")
	}

	return softDeleted, nil
}
