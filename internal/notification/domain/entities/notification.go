package entities

import (
	"time"

	"github.com/google/uuid"
)

// ─── Enum de tipo de notificación ───────────────────────────────────────────

type NotificationType string

const (
	TypeNewComment   NotificationType = "NEW_COMMENT"
	TypeAssignment   NotificationType = "ASSIGNMENT"
	TypeStatusChange NotificationType = "STATUS_CHANGE"
	TypeSystem       NotificationType = "SYSTEM"
)

// ─── Entidad principal (tabla: notifications) ───────────────────────────────

type Notification struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"user_id"`
	WorkID    *uuid.UUID       `json:"work_id,omitempty"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}
