package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/google/uuid"
)

type CreateActRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

type CreateActUseCase struct {
	repo repository.ActRepository
}

func NewCreateActUseCase(r repository.ActRepository) *CreateActUseCase {
	return &CreateActUseCase{repo: r}
}

func (uc *CreateActUseCase) Execute(ctx context.Context, req CreateActRequest) (*ActDTO, error) {
	existing, _ := uc.repo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("el nombre del acto ya está registrado")
	}

	newAct := &entities.Act{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Status:      entities.StatusActive,
	}

	if err := uc.repo.Create(ctx, newAct); err != nil {
		return nil, err
	}

	dto := ToActDTO(newAct)
	return &dto, nil
}
