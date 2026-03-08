package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresWorkRepository) AddComment(ctx context.Context, comment *entities.WorkComment) error {
	query := `
		INSERT INTO work_comments (id, work_id, user_id, message, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := repo.db.ExecContext(ctx, query,
		comment.ID, comment.WorkID, comment.UserID, comment.Message, comment.CreatedAt,
	)
	return err
}

func (repo *PostgresWorkRepository) GetCommentsByWorkID(ctx context.Context, workID uuid.UUID) ([]entities.WorkComment, error) {
	query := `
		SELECT wc.id, wc.work_id, wc.user_id, u.full_name, wc.message, wc.created_at
		FROM work_comments wc
		JOIN users u ON wc.user_id = u.id
		WHERE wc.work_id = $1
		ORDER BY wc.created_at ASC
	`
	rows, err := repo.db.QueryContext(ctx, query, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]entities.WorkComment, 0)
	for rows.Next() {
		var c entities.WorkComment
		if err := rows.Scan(&c.ID, &c.WorkID, &c.UserID, &c.UserName, &c.Message, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
