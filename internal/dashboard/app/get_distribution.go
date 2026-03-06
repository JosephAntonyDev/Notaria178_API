package app

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/domain/repository"
)

type GetDistributionUseCase struct {
	repo  repository.DashboardRepository
	cache cache.CachePort
}

func NewGetDistributionUseCase(r repository.DashboardRepository, c cache.CachePort) *GetDistributionUseCase {
	return &GetDistributionUseCase{repo: r, cache: c}
}

func (uc *GetDistributionUseCase) Execute(
	ctx context.Context,
	branchID *string,
	timeframe string,
	startDate, endDate *string,
) (*DistributionDTO, error) {

	start, end := ResolveTimeRange(timeframe, startDate, endDate)

	// ── Cache key determinista ──────────────────────────────────────────
	cacheKey := fmt.Sprintf("dashboard:dist:%s:%s:%s",
		branchKeyPart(branchID),
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

	// ── Cache hit ───────────────────────────────────────────────────────
	if uc.cache != nil {
		var cached DistributionDTO
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

	rows, err := uc.repo.GetDistribution(ctx, filters)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, r := range rows {
		total += r.Count
	}

	statuses := make([]DistributionStatusDTO, 0, len(rows))
	for _, r := range rows {
		pct := 0.0
		if total > 0 {
			pct = math.Round(float64(r.Count)/float64(total)*10000) / 100
		}
		statuses = append(statuses, DistributionStatusDTO{
			Status:     r.Status,
			Count:      r.Count,
			Percentage: pct,
		})
	}

	dto := &DistributionDTO{
		Total:    total,
		Statuses: statuses,
	}

	// ── Guardar en cache (TTL 5 min) — fire-and-forget ──────────────────
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, dto, 5*time.Minute)
	}

	return dto, nil
}
