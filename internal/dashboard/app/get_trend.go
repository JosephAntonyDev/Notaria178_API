package app

import (
	"context"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/domain/repository"
)

type GetTrendUseCase struct {
	repo  repository.DashboardRepository
	cache cache.CachePort
}

func NewGetTrendUseCase(r repository.DashboardRepository, c cache.CachePort) *GetTrendUseCase {
	return &GetTrendUseCase{repo: r, cache: c}
}

func (uc *GetTrendUseCase) Execute(
	ctx context.Context,
	userRole string,
	userID string,
	userBranchID string,
	branchID *string,
	timeframe string,
	startDate, endDate *string,
	groupBy string,
) (*TrendDTO, error) {

	// Validar y normalizar group_by
	switch groupBy {
	case "day", "week", "month":
		// ok
	default:
		groupBy = "week"
	}

	start, end := ResolveTimeRange(timeframe, startDate, endDate)

	// ── Data Scoping basado en rol ─────────────────────────────────────
	var scopedUserID *string
	switch userRole {
	case "SUPER_ADMIN":
		// Ve todo; si envía BranchID por URL se respeta.
	case "LOCAL_ADMIN":
		// Forzar aislamiento por oficina.
		branchID = &userBranchID
	case "DRAFTER", "DATA_ENTRY":
		// Solo trabajos donde es proyectista o colaborador.
		scopedUserID = &userID
	}

	// ── Cache key determinista ──────────────────────────────────────────
	cacheKey := fmt.Sprintf("dashboard:trend:%s:%s:%s:%s:%s",
		branchKeyPart(branchID),
		userKeyPart(scopedUserID),
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
		groupBy,
	)

	// ── Cache hit ───────────────────────────────────────────────────────
	if uc.cache != nil {
		var cached TrendDTO
		if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil {
			return &cached, nil
		}
	}

	// ── Cache miss → PostgreSQL ─────────────────────────────────────────
	filters := repository.DashboardFilters{
		BranchID:     branchID,
		ScopedUserID: scopedUserID,
		StartDate:    start,
		EndDate:      end,
	}

	rows, err := uc.repo.GetTrend(ctx, filters, groupBy)
	if err != nil {
		return nil, err
	}

	series := make([]TrendPointDTO, 0, len(rows))
	for _, r := range rows {
		series = append(series, TrendPointDTO{
			Period:   r.Period.Format("2006-01-02"),
			Created:  r.Created,
			Approved: r.Approved,
		})
	}

	dto := &TrendDTO{
		GroupBy: groupBy,
		Series:  series,
	}

	// ── Guardar en cache (TTL 5 min) — fire-and-forget ──────────────────
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, dto, 5*time.Minute)
	}

	return dto, nil
}
