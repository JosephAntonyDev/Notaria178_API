package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
)

func (repo *PostgresNotificationRepository) Create(ctx context.Context, notif *entities.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, work_id, type, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := repo.db.ExecContext(ctx, query,
		notif.ID, notif.UserID, notif.WorkID,
		notif.Type, notif.Message, notif.IsRead, notif.CreatedAt,
	)
	return err
}
