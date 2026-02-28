package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/ports"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
)

type UpdateProfileRequest struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone    *string `json:"phone,omitempty"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
}

type UpdateProfileUseCase struct {
	repo   repository.UserRepository
	hasher ports.PasswordHasher
}

func NewUpdateProfileUseCase(r repository.UserRepository, h ports.PasswordHasher) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{repo: r, hasher: h}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, req UpdateProfileRequest) error {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("ID inválido")
	}

	user, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil || user == nil {
		return errors.New("usuario no encontrado")
	}

	if req.Email != nil && *req.Email != user.Email {
		existing, _ := uc.repo.GetByEmail(ctx, *req.Email)
		if existing != nil {
			return errors.New("el correo electrónico ya está en uso")
		}
		user.Email = *req.Email
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if req.Password != nil && *req.Password != "" {
		hashed, err := uc.hasher.HashPassword(*req.Password)
		if err != nil {
			return errors.New("error al procesar la contraseña")
		}
		user.PasswordHash = hashed
	}

	return uc.repo.Update(ctx, user)
}