package app

import (
	"context"
	"errors"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/repository"
	"github.com/google/uuid"
)

type CreateClientUseCase struct {
	repo repository.ClientRepository
}

func NewCreateClientUseCase(r repository.ClientRepository) *CreateClientUseCase {
	return &CreateClientUseCase{repo: r}
}

func (uc *CreateClientUseCase) Execute(ctx context.Context, req CreateClientRequest) (*ClientDTO, error) {
	if req.RFC != nil && *req.RFC != "" {
		existing, _ := uc.repo.GetByRFC(ctx, *req.RFC)
		if existing != nil {
			return nil, errors.New("el RFC ya está registrado en el sistema")
		}
	}

	newClient := &entities.Client{
		ID:        uuid.New(),
		FullName:  req.FullName,
		RFC:       req.RFC,
		Phone:     req.Phone,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, newClient); err != nil {
		return nil, err
	}

	dto := ToClientDTO(newClient)
	return &dto, nil
}
