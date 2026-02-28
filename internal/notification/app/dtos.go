package app

import (
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/google/uuid"
)

// ─── DTOs ───────────────────────────────────────────────────────────────────

type NotificationDTO struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	WorkID    *uuid.UUID `json:"work_id,omitempty"`
	Type      string     `json:"type"`
	Message   string     `json:"message"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
}

// ─── Mapper ─────────────────────────────────────────────────────────────────

func ToNotificationDTO(n *entities.Notification) NotificationDTO {
	return NotificationDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		WorkID:    n.WorkID,
		Type:      string(n.Type),
		Message:   n.Message,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}
