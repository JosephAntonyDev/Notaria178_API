package entities

import (
	"time"

	"github.com/google/uuid"
)

type ActRequirement struct {
	ID        uuid.UUID `json:"id"`
	ActID     uuid.UUID `json:"act_id"`
	Name      string    `json:"name"`
	Status    ActStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
