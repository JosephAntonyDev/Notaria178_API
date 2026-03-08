package app

import (
	"context"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/domain/repository"
)

type GetKPIsUseCase struct {
	repo  repository.DashboardRepository
	cache cache.CachePort
}

func NewGetKPIsUseCase(r repository.DashboardRepository, c cache.CachePort) *GetKPIsUseCase {
	return &GetKPIsUseCase{repo: r, cache: c}
}

func (uc *GetKPIsUseCase) Execute(
	ctx context.Context,
	branchID *string,
	timeframe string,
	startDate, endDate *string,
) (*KPIsDTO, error) {

	start, end := ResolveTimeRange(timeframe, startDate, endDate)

	// ── Cache key determinista ──────────────────────────────────────────
	cacheKey := fmt.Sprintf("dashboard:kpis:%s:%s:%s",
		branchKeyPart(branchID),
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

	// ── Cache hit ───────────────────────────────────────────────────────
	if uc.cache != nil {
		var cached KPIsDTO
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

	result, err := uc.repo.GetKPIs(ctx, filters)
	if err != nil {
		return nil, err
	}

	dto := &KPIsDTO{
		Total:          result.Total,
		Pending:        result.Pending,
		InProgress:     result.InProgress,
		ReadyForReview: result.ReadyForReview,
		Approved:       result.Approved,
		Rejected:       result.Rejected,
	}

	// ── Guardar en cache (TTL 5 min) — fire-and-forget ──────────────────
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, dto, 5*time.Minute)
	}

	return dto, nil
}
