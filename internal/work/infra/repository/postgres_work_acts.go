package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresWorkRepository) AddActs(ctx context.Context, workID uuid.UUID, actIDs []uuid.UUID) error {
	query := `INSERT INTO work_acts (work_id, act_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	for _, actID := range actIDs {
		if _, err := repo.db.ExecContext(ctx, query, workID, actID); err != nil {
			return err
		}
	}
	return nil
}

func (repo *PostgresWorkRepository) RemoveAllActs(ctx context.Context, workID uuid.UUID) error {
	query := `DELETE FROM work_acts WHERE work_id = $1`
	_, err := repo.db.ExecContext(ctx, query, workID)
	return err
}

func (repo *PostgresWorkRepository) GetActsByWorkID(ctx context.Context, workID uuid.UUID) ([]entities.WorkActInfo, error) {
	query := `
		SELECT ac.id, ac.name
		FROM work_acts wa
		JOIN act_catalogs ac ON wa.act_id = ac.id
		WHERE wa.work_id = $1
		ORDER BY ac.name ASC
	`
	rows, err := repo.db.QueryContext(ctx, query, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	acts := make([]entities.WorkActInfo, 0)
	for rows.Next() {
		var act entities.WorkActInfo
		if err := rows.Scan(&act.ActID, &act.Name); err != nil {
			return nil, err
		}
		acts = append(acts, act)
	}
	return acts, nil
}
