package repository

import (
	"context"
	"strconv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

func (repo *PostgresNotificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, filters repository.NotificationFilters) ([]*entities.Notification, error) {
	baseQuery := `
		SELECT id, user_id, work_id, type, message, is_read, created_at
		FROM notifications
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argID := 2

	if filters.IsRead != nil {
		baseQuery += ` AND is_read = $` + strconv.Itoa(argID)
		args = append(args, *filters.IsRead)
		argID++
	}

	baseQuery += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifs := make([]*entities.Notification, 0)
	for rows.Next() {
		var n entities.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.WorkID, &n.Type, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifs = append(notifs, &n)
	}
	return notifs, nil
}

func (repo *PostgresNotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1 AND user_id = $2`
	_, err := repo.db.ExecContext(ctx, query, id, userID)
	return err
}

func (repo *PostgresNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`
	_, err := repo.db.ExecContext(ctx, query, userID)
	return err
}

func (repo *PostgresNotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`
	var count int
	err := repo.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
