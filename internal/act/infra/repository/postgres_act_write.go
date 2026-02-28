package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresActRepository) Create(ctx context.Context, act *entities.Act) error {
	query := `
		INSERT INTO act_catalogs (id, name, description, status)
		VALUES ($1, $2, $3, $4)
	`
	_, err := repo.db.ExecContext(ctx, query,
		act.ID, act.Name, act.Description, act.Status,
	)
	return err
}

func (repo *PostgresActRepository) Update(ctx context.Context, act *entities.Act) error {
	query := `
		UPDATE act_catalogs
		SET name = $1, description = $2
		WHERE id = $3
	`
	_, err := repo.db.ExecContext(ctx, query,
		act.Name, act.Description, act.ID,
	)
	return err
}

func (repo *PostgresActRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.ActStatus) error {
	query := `
		UPDATE act_catalogs
		SET status = $1
		WHERE id = $2
	`
	_, err := repo.db.ExecContext(ctx, query, status, id)
	return err
}
