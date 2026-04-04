package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/google/uuid"
)

type DeviceTokenRepository interface {
	// Guardar o actualizar un token FCM (upsert)
	SaveToken(ctx context.Context, token *entities.DeviceToken) error

	// Obtener todos los tokens de un usuario
	GetTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.DeviceToken, error)

	// Obtener tokens de múltiples usuarios (para notificaciones en lote)
	GetTokensByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entities.DeviceToken, error)

	// Eliminar un token específico (cuando el usuario cierra sesión)
	DeleteToken(ctx context.Context, fcmToken string) error

	// Eliminar todos los tokens de un usuario
	DeleteTokensByUserID(ctx context.Context, userID uuid.UUID) error
}
