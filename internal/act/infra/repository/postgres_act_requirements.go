package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresActRepository) AddRequirement(ctx context.Context, req *entities.ActRequirement) error {
	query := `
		INSERT INTO act_requirements (act_id, name)
		VALUES ($1, $2)
		RETURNING id, status, created_at
	`
	err := repo.db.QueryRowContext(ctx, query, req.ActID, req.Name).Scan(&req.ID, &req.Status, &req.CreatedAt)
	return err
}

func (repo *PostgresActRepository) DeleteRequirement(ctx context.Context, reqID uuid.UUID) error {
	query := `
		DELETE FROM act_requirements
		WHERE id = $1
	`
	_, err := repo.db.ExecContext(ctx, query, reqID)
	return err
}

func (repo *PostgresActRepository) DeactivateRequirement(ctx context.Context, reqID uuid.UUID) error {
	query := `UPDATE act_requirements SET status = 'INACTIVE' WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, reqID)
	return err
}

func (repo *PostgresActRepository) GetRequirements(ctx context.Context, actID uuid.UUID) ([]*entities.ActRequirement, error) {
	query := `
		SELECT id, act_id, name, status, created_at
		FROM act_requirements
		WHERE act_id = $1
		ORDER BY status ASC, created_at ASC
	`
	rows, err := repo.db.QueryContext(ctx, query, actID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requirements []*entities.ActRequirement
	for rows.Next() {
		var req entities.ActRequirement
		if err := rows.Scan(&req.ID, &req.ActID, &req.Name, &req.Status, &req.CreatedAt); err != nil {
			return nil, err
		}
		requirements = append(requirements, &req)
	}
	return requirements, nil
}
