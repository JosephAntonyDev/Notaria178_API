package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/google/uuid"
)

// SaveToken guarda o actualiza un token FCM (upsert)
func (repo *PostgresDeviceTokenRepository) SaveToken(ctx context.Context, token *entities.DeviceToken) error {
	query := `
		INSERT INTO user_device_tokens (
			id, user_id, fcm_token, device_type, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		ON CONFLICT (fcm_token)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			device_type = EXCLUDED.device_type,
			updated_at = EXCLUDED.updated_at
	`

	_, err := repo.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.FCMToken,
		token.DeviceType,
		token.CreatedAt,
		token.UpdatedAt,
	)

	return err
}

// DeleteToken elimina un token específico
func (repo *PostgresDeviceTokenRepository) DeleteToken(ctx context.Context, fcmToken string) error {
	query := `DELETE FROM user_device_tokens WHERE fcm_token = $1`
	_, err := repo.db.ExecContext(ctx, query, fcmToken)
	return err
}

// DeleteTokensByUserID elimina todos los tokens de un usuario
func (repo *PostgresDeviceTokenRepository) DeleteTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_device_tokens WHERE user_id = $1`
	_, err := repo.db.ExecContext(ctx, query, userID)
	return err
}
