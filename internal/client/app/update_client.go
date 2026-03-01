package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/repository"
	"github.com/google/uuid"
)

type UpdateClientUseCase struct {
	repo repository.ClientRepository
}

func NewUpdateClientUseCase(r repository.ClientRepository) *UpdateClientUseCase {
	return &UpdateClientUseCase{repo: r}
}

func (uc *UpdateClientUseCase) Execute(ctx context.Context, clientID string, req UpdateClientRequest) (*ClientDTO, error) {
	parsedID, err := uuid.Parse(clientID)
	if err != nil {
		return nil, errors.New("ID de cliente inválido")
	}

	client, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("cliente no encontrado")
	}

	if req.RFC != nil && *req.RFC != "" {
		currentRFC := ""
		if client.RFC != nil {
			currentRFC = *client.RFC
		}
		if *req.RFC != currentRFC {
			existing, _ := uc.repo.GetByRFC(ctx, *req.RFC)
			if existing != nil {
				return nil, errors.New("el RFC ya está registrado en el sistema")
			}
		}
		client.RFC = req.RFC
	}

	if req.FullName != nil {
		client.FullName = *req.FullName
	}
	if req.Phone != nil {
		client.Phone = req.Phone
	}
	if req.Email != nil {
		client.Email = req.Email
	}

	if err := uc.repo.Update(ctx, client); err != nil {
		return nil, err
	}

	dto := ToClientDTO(client)
	return &dto, nil
}
