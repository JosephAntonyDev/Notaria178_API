package entities

import (
	"time"

	"github.com/google/uuid"
)

// ─── Entidad de comentarios del expediente (tabla: work_comments) ───────────

type WorkComment struct {
	ID        uuid.UUID `json:"id"`
	WorkID    uuid.UUID `json:"work_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"` // se llena via JOIN en lectura
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
