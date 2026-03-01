package entities

import "github.com/google/uuid"

// ─── Value Object para JOIN work_acts ↔ act_catalogs (tabla: work_acts) ─────

type WorkActInfo struct {
	ActID uuid.UUID `json:"act_id"`
	Name  string    `json:"name"`
}
