package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// GetTokensByUserID obtiene todos los tokens de un usuario
func (repo *PostgresDeviceTokenRepository) GetTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.DeviceToken, error) {
	query := `
		SELECT id, user_id, fcm_token, device_type, created_at, updated_at
		FROM user_device_tokens
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entities.DeviceToken
	for rows.Next() {
		token := &entities.DeviceToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.FCMToken,
			&token.DeviceType,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetTokensByUserIDs obtiene tokens de múltiples usuarios (para notificaciones en lote)
func (repo *PostgresDeviceTokenRepository) GetTokensByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*entities.DeviceToken, error) {
	if len(userIDs) == 0 {
		return []*entities.DeviceToken{}, nil
	}

	query := `
		SELECT id, user_id, fcm_token, device_type, created_at, updated_at
		FROM user_device_tokens
		WHERE user_id = ANY($1)
		ORDER BY user_id, updated_at DESC
	`

	rows, err := repo.db.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entities.DeviceToken
	for rows.Next() {
		token := &entities.DeviceToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.FCMToken,
			&token.DeviceType,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}
