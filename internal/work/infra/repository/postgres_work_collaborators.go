package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresWorkRepository) AddCollaborator(ctx context.Context, workID uuid.UUID, userID uuid.UUID) error {
	query := `INSERT INTO work_collaborators (work_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := repo.db.ExecContext(ctx, query, workID, userID)
	return err
}

func (repo *PostgresWorkRepository) RemoveCollaborator(ctx context.Context, workID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM work_collaborators WHERE work_id = $1 AND user_id = $2`
	_, err := repo.db.ExecContext(ctx, query, workID, userID)
	return err
}

func (repo *PostgresWorkRepository) GetCollaborators(ctx context.Context, workID uuid.UUID) ([]entities.WorkCollaboratorInfo, error) {
	query := `
		SELECT u.id, u.full_name
		FROM work_collaborators wc
		JOIN users u ON wc.user_id = u.id
		WHERE wc.work_id = $1
		ORDER BY u.full_name ASC
	`
	rows, err := repo.db.QueryContext(ctx, query, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collabs := make([]entities.WorkCollaboratorInfo, 0)
	for rows.Next() {
		var c entities.WorkCollaboratorInfo
		if err := rows.Scan(&c.UserID, &c.FullName); err != nil {
			return nil, err
		}
		collabs = append(collabs, c)
	}
	return collabs, nil
}

func (repo *PostgresWorkRepository) IsCollaborator(ctx context.Context, workID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM work_collaborators WHERE work_id = $1 AND user_id = $2)`
	var exists bool
	err := repo.db.QueryRowContext(ctx, query, workID, userID).Scan(&exists)
	return exists, err
}
