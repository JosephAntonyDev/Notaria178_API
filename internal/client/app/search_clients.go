package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/repository"
)

type SearchClientsUseCase struct {
	repo repository.ClientRepository
}

func NewSearchClientsUseCase(r repository.ClientRepository) *SearchClientsUseCase {
	return &SearchClientsUseCase{repo: r}
}

func (uc *SearchClientsUseCase) Execute(ctx context.Context, filters repository.ClientFilters) ([]ClientDTO, error) {
	clients, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	clientsDTO := make([]ClientDTO, 0, len(clients))
	for _, c := range clients {
		clientsDTO = append(clientsDTO, ToClientDTO(c))
	}

	return clientsDTO, nil
}
