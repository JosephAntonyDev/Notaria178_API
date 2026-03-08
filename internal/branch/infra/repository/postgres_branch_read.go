package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/repository"
	"github.com/google/uuid"
)

func (repo *PostgresBranchRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Branch, error) {
	query := `
		SELECT id, name, address, created_at
		FROM branches
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	var branch entities.Branch
	err := row.Scan(&branch.ID, &branch.Name, &branch.Address, &branch.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &branch, nil
}

func (repo *PostgresBranchRepository) GetByName(ctx context.Context, name string) (*entities.Branch, error) {
	query := `
		SELECT id, name, address, created_at
		FROM branches
		WHERE name = $1
	`
	row := repo.db.QueryRowContext(ctx, query, name)
	var branch entities.Branch
	err := row.Scan(&branch.ID, &branch.Name, &branch.Address, &branch.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &branch, nil
}

func (repo *PostgresBranchRepository) List(ctx context.Context, filters repository.BranchFilters) ([]*entities.Branch, error) {
	baseQuery := `
		SELECT id, name, address, created_at
		FROM branches
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.Search != nil && *filters.Search != "" {
		baseQuery += ` AND (name ILIKE $` + strconv.Itoa(argID) + ` OR address ILIKE $` + strconv.Itoa(argID) + `)`
		args = append(args, "%"+*filters.Search+"%")
		argID++
	}

	baseQuery += ` ORDER BY name ASC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []*entities.Branch
	for rows.Next() {
		var branch entities.Branch
		err := rows.Scan(&branch.ID, &branch.Name, &branch.Address, &branch.CreatedAt)
		if err != nil {
			return nil, err
		}
		branches = append(branches, &branch)
	}
	return branches, nil
}
