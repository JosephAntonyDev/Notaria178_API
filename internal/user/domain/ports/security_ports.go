package ports

import (
	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword string, providedPassword string) error
}

type TokenManager interface {
	GenerateToken(userID uuid.UUID, role entities.UserRole, branchID *uuid.UUID) (string, error)
	ValidateToken(token string) (bool, map[string]interface{}, error)
}