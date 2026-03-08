package app

import (
	"context"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/domain/repository"
)

type GetTopActsUseCase struct {
	repo  repository.DashboardRepository
	cache cache.CachePort
}

func NewGetTopActsUseCase(r repository.DashboardRepository, c cache.CachePort) *GetTopActsUseCase {
	return &GetTopActsUseCase{repo: r, cache: c}
}

func (uc *GetTopActsUseCase) Execute(
	ctx context.Context,
	branchID *string,
	timeframe string,
	startDate, endDate *string,
	limit int,
) (*TopActsDTO, error) {

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	start, end := ResolveTimeRange(timeframe, startDate, endDate)

	// ── Cache key determinista ──────────────────────────────────────────
	cacheKey := fmt.Sprintf("dashboard:topacts:%s:%s:%s:%d",
		branchKeyPart(branchID),
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
		limit,
	)

	// ── Cache hit ───────────────────────────────────────────────────────
	if uc.cache != nil {
		var cached TopActsDTO
		if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil {
			return &cached, nil
		}
	}

	// ── Cache miss → PostgreSQL ─────────────────────────────────────────
	filters := repository.DashboardFilters{
		BranchID:  branchID,
		StartDate: start,
		EndDate:   end,
	}

	rows, err := uc.repo.GetTopActs(ctx, filters, limit)
	if err != nil {
		return nil, err
	}

	data := make([]TopActDTO, 0, len(rows))
	for _, r := range rows {
		data = append(data, TopActDTO{
			ActID: r.ActID,
			Name:  r.Name,
			Count: r.Count,
		})
	}

	dto := &TopActsDTO{Data: data}

	// ── Guardar en cache (TTL 5 min) — fire-and-forget ──────────────────
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, dto, 5*time.Minute)
	}

	return dto, nil
}
