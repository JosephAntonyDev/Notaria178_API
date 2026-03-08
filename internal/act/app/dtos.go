package app

import (
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/google/uuid"
)

type ActDTO struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	Category          string    `json:"category"`
	Status            string    `json:"status"`
	RequirementsCount int       `json:"requirements_count"`
	WorksCount        int       `json:"works_count"`
}

func ToActDTO(act *entities.Act) ActDTO {
	return ActDTO{
		ID:                act.ID,
		Name:              act.Name,
		Description:       act.Description,
		Category:          act.Category,
		Status:            string(act.Status),
		RequirementsCount: act.RequirementsCount,
		WorksCount:        act.WorksCount,
	}
}
