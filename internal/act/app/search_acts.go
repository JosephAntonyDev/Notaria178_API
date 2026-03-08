package app

import (
	"context"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
)

type SearchActsUseCase struct {
	repo  repository.ActRepository
	cache cache.CachePort
}

func NewSearchActsUseCase(r repository.ActRepository, c cache.CachePort) *SearchActsUseCase {
	return &SearchActsUseCase{repo: r, cache: c}
}

func (uc *SearchActsUseCase) Execute(ctx context.Context, filters repository.ActFilters) ([]ActDTO, error) {
	// Construir cache key determinista
	cacheKey := fmt.Sprintf("acts:search:%d:%d:%v:%v",
		filters.Limit, filters.Offset,
		ptrVal(filters.Search), ptrVal(filters.Status),
	)

	// Cache Hit
	if uc.cache != nil {
		var cached []ActDTO
		if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil {
			return cached, nil
		}
		// Cache miss o error de Redis: continuar a BD
	}

	acts, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	actsDTO := make([]ActDTO, 0, len(acts))
	for _, a := range acts {
		actsDTO = append(actsDTO, ToActDTO(a))
	}

	// Guardar en cache (TTL 24h) — fire-and-forget
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, actsDTO, 24*time.Hour)
	}

	return actsDTO, nil
}

// ptrVal devuelve el valor de un *string o "nil" si es nulo.
func ptrVal(s *string) string {
	if s == nil {
		return "nil"
	}
	return *s
}
