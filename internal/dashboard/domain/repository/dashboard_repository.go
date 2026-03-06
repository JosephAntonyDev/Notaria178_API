package repository

import (
	"context"
	"time"
)

// ─── Filtros comunes para todos los endpoints del dashboard ─────────────────

type DashboardFilters struct {
	BranchID  *string
	StartDate time.Time
	EndDate   time.Time
}

// ─── Resultados de KPIs ─────────────────────────────────────────────────────

type KPIsResult struct {
	Total          int
	Pending        int
	InProgress     int
	ReadyForReview int
	Approved       int
	Rejected       int
}

// ─── Resultados de Tendencia ────────────────────────────────────────────────

type TrendRow struct {
	Period   time.Time
	Created  int
	Approved int
}

// ─── Resultados de Distribución ─────────────────────────────────────────────

type DistributionRow struct {
	Status string
	Count  int
}

// ─── Filtros para Actividad Reciente ────────────────────────────────────────

type ActivityFilters struct {
	BranchID  *string
	UserID    *string
	EntityID  *string
	StartDate time.Time
	EndDate   time.Time
	Limit     int
	Offset    int
}

// ─── Resultados de Actividad Reciente ───────────────────────────────────────

type ActivityRow struct {
	ID          string
	UserID      *string
	UserName    *string
	Action      string
	Entity      string
	EntityID    string
	JSONDetails []byte
	CreatedAt   time.Time
}

// ─── Resultados de Top Drafters ─────────────────────────────────────────────

type TopDrafterRow struct {
	UserID    string
	FullName  string
	Role      string
	WorkCount int
}

// ─── Resultados de Top Acts ─────────────────────────────────────────────────

type TopActRow struct {
	ActID string
	Name  string
	Count int
}

// ─── Puerto de salida (interfaz del repositorio) ────────────────────────────

type DashboardRepository interface {
	GetKPIs(ctx context.Context, filters DashboardFilters) (*KPIsResult, error)
	GetTrend(ctx context.Context, filters DashboardFilters, groupBy string) ([]TrendRow, error)
	GetDistribution(ctx context.Context, filters DashboardFilters) ([]DistributionRow, error)
	GetRecentActivity(ctx context.Context, filters ActivityFilters) ([]ActivityRow, int, error)
	GetTopDrafters(ctx context.Context, filters DashboardFilters, limit int) ([]TopDrafterRow, error)
	GetTopActs(ctx context.Context, filters DashboardFilters, limit int) ([]TopActRow, error)
}
