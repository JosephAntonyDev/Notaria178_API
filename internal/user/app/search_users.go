package app

import (
	"context"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

type SearchUsersUseCase struct {
	repo repository.UserRepository
}

func NewSearchUsersUseCase(r repository.UserRepository) *SearchUsersUseCase {
	return &SearchUsersUseCase{repo: r}
}

// Ahora recibe repository.UserFilters
func (uc *SearchUsersUseCase) Execute(ctx context.Context, filters entities.UserFilters) ([]UserPublicDTO, error) {
	users, err := uc.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	usersDTO := make([]UserPublicDTO, 0)
	for _, u := range users {
		usersDTO = append(usersDTO, ToUserPublicDTO(u))
	}

	return usersDTO, nil
}