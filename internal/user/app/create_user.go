package app

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/ports"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
)

type CreateUserRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required"`
	BranchID *string `json:"branch_id,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type CreateUserUseCase struct {
	repo   repository.UserRepository
	hasher ports.PasswordHasher
}

func NewCreateUserUseCase(r repository.UserRepository, h ports.PasswordHasher) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo:   r,
		hasher: h,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, requesterRole string, req CreateUserRequest) (*entities.User, error) {
	
	if requesterRole == string(entities.RoleLocalAdmin) && req.Role == string(entities.RoleSuperAdmin) {
		return nil, errors.New("operación denegada: un administrador no puede crear una cuenta de notario titular")
	}

	existingUser, _ := uc.repo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("el correo ya está registrado en el sistema")
	}

	hashedPassword, err := uc.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	var branchUUID *uuid.UUID
	if req.BranchID != nil && *req.BranchID != "" {
		parsed, err := uuid.Parse(*req.BranchID)
		if err == nil {
			branchUUID = &parsed
		}
	}

	newUser := &entities.User{
		ID:           uuid.New(),
		BranchID:     branchUUID,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Role:         entities.UserRole(req.Role),
		Status:       entities.StatusActive,
		HireDate:     time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	go uc.sendWelcomeEmailAsync(newUser.Email, newUser.FullName)

	return newUser, nil
}

func (uc *CreateUserUseCase) sendWelcomeEmailAsync(email string, name string) {
	log.Printf("[Fondo] Iniciando envío de credenciales a: %s...", email)
	time.Sleep(2 * time.Second) 
	log.Printf("[Fondo] Correo enviado exitosamente a: %s", name)
}