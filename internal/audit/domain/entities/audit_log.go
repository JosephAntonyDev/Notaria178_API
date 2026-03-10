package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ─── Entidad principal (tabla: audit_logs) ──────────────────────────────────

type AuditLog struct {
	ID          uuid.UUID       `json:"id"`
	UserID      *uuid.UUID      `json:"user_id,omitempty"`
	UserName    *string         `json:"user_name,omitempty"`
	Action      string          `json:"action"`
	Entity      string          `json:"entity"`
	EntityID    uuid.UUID       `json:"entity_id"`
	JSONDetails json.RawMessage `json:"json_details,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ─── Métricas Agrupadas ─────────────────────────────────────────────────────

type ActionMetricItem struct {
	Action string `json:"action"`
	Count  int    `json:"count"`
}
