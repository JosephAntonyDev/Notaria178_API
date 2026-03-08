package app

import (
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/entities"
	"github.com/google/uuid"
)

type BranchDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   *string   `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func ToBranchDTO(branch *entities.Branch) BranchDTO {
	return BranchDTO{
		ID:        branch.ID,
		Name:      branch.Name,
		Address:   branch.Address,
		CreatedAt: branch.CreatedAt,
	}
}

type CreateBranchRequest struct {
	Name    string  `json:"name" binding:"required"`
	Address *string `json:"address,omitempty"`
}

type UpdateBranchRequest struct {
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
}
