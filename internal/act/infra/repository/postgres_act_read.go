package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/google/uuid"
)

func (repo *PostgresActRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Act, error) {
	query := `
		SELECT id, name, description, status
		FROM act_catalogs
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	var act entities.Act
	err := row.Scan(&act.ID, &act.Name, &act.Description, &act.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &act, nil
}

func (repo *PostgresActRepository) GetByName(ctx context.Context, name string) (*entities.Act, error) {
	query := `
		SELECT id, name, description, status
		FROM act_catalogs
		WHERE name = $1
	`
	row := repo.db.QueryRowContext(ctx, query, name)
	var act entities.Act
	err := row.Scan(&act.ID, &act.Name, &act.Description, &act.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &act, nil
}

func (repo *PostgresActRepository) List(ctx context.Context, filters repository.ActFilters) ([]*entities.Act, error) {
	baseQuery := `
		SELECT id, name, description, status
		FROM act_catalogs
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.Search != nil && *filters.Search != "" {
		baseQuery += ` AND name ILIKE $` + strconv.Itoa(argID)
		args = append(args, "%"+*filters.Search+"%")
		argID++
	}
	if filters.Status != nil && *filters.Status != "" {
		baseQuery += ` AND status = $` + strconv.Itoa(argID)
		args = append(args, *filters.Status)
		argID++
	}

	baseQuery += ` ORDER BY name ASC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var acts []*entities.Act
	for rows.Next() {
		var act entities.Act
		err := rows.Scan(&act.ID, &act.Name, &act.Description, &act.Status)
		if err != nil {
			return nil, err
		}
		acts = append(acts, &act)
	}
	return acts, nil
}
