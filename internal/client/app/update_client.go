package app

import (
	"context"
	"errors"
	"time"

	clientEntities "github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/entities"
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

	// Validar RFC único si se envía uno nuevo
	newRFC := ""
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
		newRFC = *req.RFC
	}

	// Verificar si el cliente tiene trabajos APPROVED (Copy-on-Write)
	approvedCount, _ := uc.repo.CountWorksWithClientInStatus(ctx, parsedID, "APPROVED")

	if approvedCount > 0 {
		// COPY-ON-WRITE: crear nuevo registro con datos editados
		newClient := &clientEntities.Client{
			ID:        uuid.New(),
			FullName:  client.FullName,
			RFC:       client.RFC,
			Phone:     client.Phone,
			Email:     client.Email,
			CreatedAt: time.Now(),
		}

		// Aplicar cambios al nuevo registro
		if req.FullName != nil {
			newClient.FullName = *req.FullName
		}
		if req.RFC != nil {
			if newRFC != "" {
				newClient.RFC = &newRFC
			} else {
				newClient.RFC = req.RFC
			}
		}
		if req.Phone != nil {
			newClient.Phone = req.Phone
		}
		if req.Email != nil {
			newClient.Email = req.Email
		}

		if err := uc.repo.Create(ctx, newClient); err != nil {
			return nil, errors.New("error al crear copia del cliente")
		}

		// Mover trabajos no-APPROVED al nuevo cliente
		if err := uc.repo.UpdatePendingWorksClientID(ctx, parsedID, newClient.ID); err != nil {
			return nil, errors.New("error al actualizar trabajos pendientes")
		}

		dto := ToClientDTO(newClient)
		return &dto, nil
	}

	// Sin trabajos aprobados: actualización directa
	if req.FullName != nil {
		client.FullName = *req.FullName
	}
	if req.RFC != nil {
		client.RFC = req.RFC
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
