package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/entities"
)

// ─── Filtros dinámicos ──────────────────────────────────────────────────────

type AuditFilters struct {
	Limit     int
	Offset    int
	UserID    *string
	Action    *string
	Entity    *string
	EntityID  *string
	StartDate *string
	EndDate   *string
}

// ─── Puerto de salida (interfaz del repositorio) ────────────────────────────

type AuditRepository interface {
	Create(ctx context.Context, log *entities.AuditLog) error
	List(ctx context.Context, filters AuditFilters) ([]*entities.AuditLog, error)
}
