package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/entities"
)

func (repo *PostgresBranchRepository) Create(ctx context.Context, branch *entities.Branch) error {
	query := `
		INSERT INTO branches (id, name, address, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return repo.db.QueryRowContext(ctx, query,
		branch.ID, branch.Name, branch.Address, branch.CreatedAt,
	).Scan(&branch.ID)
}

func (repo *PostgresBranchRepository) Update(ctx context.Context, branch *entities.Branch) error {
	query := `
		UPDATE branches
		SET name = $1, address = $2
		WHERE id = $3
	`
	_, err := repo.db.ExecContext(ctx, query,
		branch.Name, branch.Address, branch.ID,
	)
	return err
}
