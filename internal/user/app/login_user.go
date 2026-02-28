package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/ports"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
)

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string         `json:"token"`
	User  *entities.User `json:"user"`
}

type LoginUserUseCase struct {
	repo         repository.UserRepository
	hasher       ports.PasswordHasher
	tokenManager ports.TokenManager
}

func NewLoginUserUseCase(r repository.UserRepository, h ports.PasswordHasher, t ports.TokenManager) *LoginUserUseCase {
	return &LoginUserUseCase{
		repo:         r,
		hasher:       h,
		tokenManager: t,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, req LoginUserRequest) (*LoginResponse, error) {
	user, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("credenciales inválidas")
	}

	if user.Status != entities.StatusActive {
		return nil, errors.New("esta cuenta ha sido desactivada, contacta al administrador")
	}

	err = uc.hasher.ComparePasswords(user.PasswordHash, req.Password)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	token, err := uc.tokenManager.GenerateToken(user.ID, user.Role, user.BranchID)
	if err != nil {
		return nil, errors.New("error interno al generar el token de acceso")
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}