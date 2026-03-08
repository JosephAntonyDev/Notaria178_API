package app

import (
	"encoding/json"
	"time"
)

// ─── Timeframe helpers ──────────────────────────────────────────────────────

// ResolveTimeRange convierte el parámetro "timeframe" o las fechas manuales
// en un par (start, end) concretos. Las fechas manuales tienen prioridad.
func ResolveTimeRange(timeframe string, startDate, endDate *string) (time.Time, time.Time) {
	// Si se proporcionan fechas manuales, prioridad absoluta.
	if startDate != nil && *startDate != "" && endDate != nil && *endDate != "" {
		s, errS := time.Parse("2006-01-02", *startDate)
		e, errE := time.Parse("2006-01-02", *endDate)
		if errS == nil && errE == nil {
			// end_date se incluye completo (hasta las 23:59:59)
			return s, e.AddDate(0, 0, 1)
		}
	}

	now := time.Now()

	switch timeframe {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return start, start.AddDate(0, 0, 1)
	case "week":
		return now.AddDate(0, 0, -7), now
	case "3months":
		return now.AddDate(0, -3, 0), now
	case "6months":
		return now.AddDate(0, -6, 0), now
	case "9months":
		return now.AddDate(0, -9, 0), now
	case "year":
		return now.AddDate(-1, 0, 0), now
	case "all":
		// Sin límite inferior: usamos una fecha mínima razonable.
		return time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location()), now
	case "month":
		// Default: mes actual completo.
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, now
	default:
		// Si no se envió nada, default = mes actual.
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, now
	}
}

// ─── DTOs de respuesta: KPIs ────────────────────────────────────────────────

type KPIsDTO struct {
	Total          int `json:"total"`
	Pending        int `json:"pending"`
	InProgress     int `json:"in_progress"`
	ReadyForReview int `json:"ready_for_review"`
	Approved       int `json:"approved"`
	Rejected       int `json:"rejected"`
}

// ─── DTOs de respuesta: Trend ───────────────────────────────────────────────

type TrendPointDTO struct {
	Period   string `json:"period"`
	Created  int    `json:"created"`
	Approved int    `json:"approved"`
}

type TrendDTO struct {
	GroupBy string          `json:"group_by"`
	Series  []TrendPointDTO `json:"series"`
}

// ─── DTOs de respuesta: Distribution ────────────────────────────────────────

type DistributionStatusDTO struct {
	Status     string  `json:"status"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type DistributionDTO struct {
	Total    int                     `json:"total"`
	Statuses []DistributionStatusDTO `json:"statuses"`
}

// ─── DTOs de respuesta: Activity ────────────────────────────────────────────

type ActivityItemDTO struct {
	ID          string          `json:"id"`
	UserID      *string         `json:"user_id,omitempty"`
	UserName    *string         `json:"user_name,omitempty"`
	Action      string          `json:"action"`
	Entity      string          `json:"entity"`
	EntityID    string          `json:"entity_id"`
	JSONDetails json.RawMessage `json:"json_details,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ActivityDTO struct {
	Total int               `json:"total"`
	Data  []ActivityItemDTO `json:"data"`
}

// ─── DTOs de respuesta: Top Drafters ────────────────────────────────────────

type TopDrafterDTO struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	WorkCount int    `json:"work_count"`
}

type TopDraftersDTO struct {
	Data []TopDrafterDTO `json:"data"`
}

// ─── DTOs de respuesta: Top Acts ────────────────────────────────────────────

type TopActDTO struct {
	ActID string `json:"act_id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TopActsDTO struct {
	Data []TopActDTO `json:"data"`
}

// ─── Helper para cache keys ─────────────────────────────────────────────────

func branchKeyPart(branchID *string) string {
	if branchID == nil || *branchID == "" {
		return "all"
	}
	return *branchID
}

func ptrKeyPart(s *string) string {
	if s == nil || *s == "" {
		return "_"
	}
	return *s
}
