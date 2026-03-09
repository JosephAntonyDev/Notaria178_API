package repository

import (
	"context"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresWorkRepository) AddWorkRequirement(ctx context.Context, workID uuid.UUID, name string) (*entities.WorkRequirement, error) {
	req := &entities.WorkRequirement{
		ID:        uuid.New(),
		WorkID:    workID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	query := `INSERT INTO work_requirements (id, work_id, name, created_at) VALUES ($1, $2, $3, $4)`
	_, err := repo.db.ExecContext(ctx, query, req.ID, req.WorkID, req.Name, req.CreatedAt)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (repo *PostgresWorkRepository) GetWorkRequirements(ctx context.Context, workID uuid.UUID) ([]entities.WorkRequirement, error) {
	query := `
		SELECT id, work_id, name, document_id, created_at
		FROM work_requirements
		WHERE work_id = $1
		ORDER BY created_at ASC
	`
	rows, err := repo.db.QueryContext(ctx, query, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reqs := make([]entities.WorkRequirement, 0)
	for rows.Next() {
		var r entities.WorkRequirement
		if err := rows.Scan(&r.ID, &r.WorkID, &r.Name, &r.DocumentID, &r.CreatedAt); err != nil {
			return nil, err
		}
		reqs = append(reqs, r)
	}
	return reqs, nil
}

func (repo *PostgresWorkRepository) DeleteWorkRequirement(ctx context.Context, reqID uuid.UUID) error {
	query := `DELETE FROM work_requirements WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, reqID)
	return err
}

func (repo *PostgresWorkRepository) GetWorkRequirementByID(ctx context.Context, reqID uuid.UUID) (*entities.WorkRequirement, error) {
	query := `SELECT id, work_id, name, document_id, created_at FROM work_requirements WHERE id = $1`
	row := repo.db.QueryRowContext(ctx, query, reqID)
	var r entities.WorkRequirement
	if err := row.Scan(&r.ID, &r.WorkID, &r.Name, &r.DocumentID, &r.CreatedAt); err != nil {
		return nil, err
	}
	return &r, nil
}
