package app

import (
	"encoding/json"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/entities"
	"github.com/google/uuid"
)

// ─── DTOs ───────────────────────────────────────────────────────────────────

type AuditLogDTO struct {
	ID          uuid.UUID       `json:"id"`
	UserID      *uuid.UUID      `json:"user_id,omitempty"`
	Action      string          `json:"action"`
	Entity      string          `json:"entity"`
	EntityID    uuid.UUID       `json:"entity_id"`
	JSONDetails json.RawMessage `json:"json_details,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ─── Mapper ─────────────────────────────────────────────────────────────────

func ToAuditLogDTO(l *entities.AuditLog) AuditLogDTO {
	return AuditLogDTO{
		ID:          l.ID,
		UserID:      l.UserID,
		Action:      l.Action,
		Entity:      l.Entity,
		EntityID:    l.EntityID,
		JSONDetails: l.JSONDetails,
		CreatedAt:   l.CreatedAt,
	}
}
