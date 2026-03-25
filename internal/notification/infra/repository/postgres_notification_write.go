package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
)

func (repo *PostgresNotificationRepository) Create(ctx context.Context, notif *entities.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, work_id, type, title, body, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := repo.db.ExecContext(ctx, query,
		notif.ID, notif.UserID, notif.WorkID,
		notif.Type, notif.Title, notif.Body, notif.Message, notif.IsRead, notif.CreatedAt,
	)
	return err
}

// CreateBatch crea múltiples notificaciones en una sola transacción
func (repo *PostgresNotificationRepository) CreateBatch(ctx context.Context, notifs []*entities.Notification) error {
	if len(notifs) == 0 {
		return nil
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO notifications (id, user_id, work_id, type, title, body, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, notif := range notifs {
		_, err := stmt.ExecContext(ctx,
			notif.ID, notif.UserID, notif.WorkID,
			notif.Type, notif.Title, notif.Body, notif.Message, notif.IsRead, notif.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
