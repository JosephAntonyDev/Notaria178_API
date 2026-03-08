package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/google/uuid"
)

type ActRequirementDTO struct {
	ID        uuid.UUID `json:"id"`
	ActID     uuid.UUID `json:"act_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
}

func ToActRequirementDTO(req *entities.ActRequirement) ActRequirementDTO {
	return ActRequirementDTO{
		ID:        req.ID,
		ActID:     req.ActID,
		Name:      req.Name,
		Status:    string(req.Status),
		CreatedAt: req.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type GetRequirementsUseCase struct {
	repo repository.ActRepository
}

func NewGetRequirementsUseCase(r repository.ActRepository) *GetRequirementsUseCase {
	return &GetRequirementsUseCase{repo: r}
}

func (uc *GetRequirementsUseCase) Execute(ctx context.Context, actID uuid.UUID) ([]ActRequirementDTO, error) {
	reqs, err := uc.repo.GetRequirements(ctx, actID)
	if err != nil {
		return nil, err
	}

	dtos := make([]ActRequirementDTO, 0, len(reqs))
	for _, req := range reqs {
		dtos = append(dtos, ToActRequirementDTO(req))
	}
	return dtos, nil
}
