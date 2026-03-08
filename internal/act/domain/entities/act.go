package entities

import "github.com/google/uuid"

type ActStatus string

const (
	StatusActive   ActStatus = "ACTIVE"
	StatusInactive ActStatus = "INACTIVE"
)

type Act struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	Category          string    `json:"category"`
	Status            ActStatus `json:"status"`
	RequirementsCount int       `json:"requirements_count"`
	WorksCount        int       `json:"works_count"`
}
