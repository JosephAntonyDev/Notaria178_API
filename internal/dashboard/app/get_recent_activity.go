package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/domain/repository"
)

type GetRecentActivityUseCase struct {
	repo  repository.DashboardRepository
	cache cache.CachePort
}

func NewGetRecentActivityUseCase(r repository.DashboardRepository, c cache.CachePort) *GetRecentActivityUseCase {
	return &GetRecentActivityUseCase{repo: r, cache: c}
}

func (uc *GetRecentActivityUseCase) Execute(
	ctx context.Context,
	branchID *string,
	userID *string,
	entityID *string,
	timeframe string,
	startDate, endDate *string,
	limit, offset int,
) (*ActivityDTO, error) {

	start, end := ResolveTimeRange(timeframe, startDate, endDate)

	// ── Cache key determinista ──────────────────────────────────────────
	cacheKey := fmt.Sprintf("dashboard:activity:%s:%s:%s:%s:%s:%d:%d",
		branchKeyPart(branchID),
		ptrKeyPart(userID),
		ptrKeyPart(entityID),
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
		limit, offset,
	)

	// ── Cache hit ───────────────────────────────────────────────────────
	if uc.cache != nil {
		var cached ActivityDTO
		if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil {
			return &cached, nil
		}
	}

	// ── Cache miss → PostgreSQL ─────────────────────────────────────────
	filters := repository.ActivityFilters{
		BranchID:  branchID,
		UserID:    userID,
		EntityID:  entityID,
		StartDate: start,
		EndDate:   end,
		Limit:     limit,
		Offset:    offset,
	}

	rows, total, err := uc.repo.GetRecentActivity(ctx, filters)
	if err != nil {
		return nil, err
	}

	items := make([]ActivityItemDTO, 0, len(rows))
	for _, r := range rows {
		item := ActivityItemDTO{
			ID:        r.ID,
			UserID:    r.UserID,
			UserName:  r.UserName,
			Action:    r.Action,
			Entity:    r.Entity,
			EntityID:  r.EntityID,
			CreatedAt: r.CreatedAt,
		}
		if len(r.JSONDetails) > 0 {
			item.JSONDetails = json.RawMessage(r.JSONDetails)
		}
		items = append(items, item)
	}

	dto := &ActivityDTO{
		Total: total,
		Data:  items,
	}

	// ── Guardar en cache (TTL 3 min) — fire-and-forget ──────────────────
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, dto, 3*time.Minute)
	}

	return dto, nil
}
