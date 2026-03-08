package repository

import (
	"context"
	"log"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresActRepository) Create(ctx context.Context, act *entities.Act) error {
	query := `
		INSERT INTO act_catalogs (id, name, description, category, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	log.Printf("[DEBUG] Create Act: ID=%s Name=%s Category=%s Status=%s Desc=%v",
		act.ID, act.Name, act.Category, act.Status, act.Description)
	_, err := repo.db.ExecContext(ctx, query,
		act.ID, act.Name, act.Description, act.Category, act.Status,
	)
	if err != nil {
		log.Printf("[ERROR] Create Act ExecContext failed: %v", err)
	}
	return err
}

func (repo *PostgresActRepository) Update(ctx context.Context, act *entities.Act) error {
	query := `
		UPDATE act_catalogs
		SET name = $1, description = $2, category = $3
		WHERE id = $4
	`
	_, err := repo.db.ExecContext(ctx, query,
		act.Name, act.Description, act.Category, act.ID,
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

func (repo *PostgresActRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM act_catalogs WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}
