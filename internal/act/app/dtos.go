package app

import (
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/google/uuid"
)

type ActDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
}

func ToActDTO(act *entities.Act) ActDTO {
	return ActDTO{
		ID:          act.ID,
		Name:        act.Name,
		Description: act.Description,
		Status:      string(act.Status),
	}
}
