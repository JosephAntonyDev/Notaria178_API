package repository

import (
	"context"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresWorkRepository) Create(ctx context.Context, work *entities.Work) error {
	query := `
		INSERT INTO works (id, branch_id, client_id, main_drafter_id, folio, status, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := repo.db.ExecContext(ctx, query,
		work.ID, work.BranchID, work.ClientID, work.MainDrafterID,
		work.Folio, work.Status, work.Deadline, work.CreatedAt, work.UpdatedAt,
	)
	return err
}

func (repo *PostgresWorkRepository) Update(ctx context.Context, work *entities.Work) error {
	query := `
		UPDATE works
		SET folio = $1, deadline = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := repo.db.ExecContext(ctx, query,
		work.Folio, work.Deadline, time.Now(), work.ID,
	)
	return err
}

func (repo *PostgresWorkRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.WorkStatus) error {
	query := `
		UPDATE works
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := repo.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}
