package app

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
)

type SearchUsersUseCase struct {
	repo repository.UserRepository
}

func NewSearchUsersUseCase(r repository.UserRepository) *SearchUsersUseCase {
	return &SearchUsersUseCase{
		repo: r,
	}
}

func (uc *SearchUsersUseCase) Execute(ctx context.Context, limit int, offset int) ([]UserPublicDTO, error) {
	users, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	usersDTO := make([]UserPublicDTO, 0)
	
	for _, u := range users {
		usersDTO = append(usersDTO, ToUserPublicDTO(u))
	}

	return usersDTO, nil
}