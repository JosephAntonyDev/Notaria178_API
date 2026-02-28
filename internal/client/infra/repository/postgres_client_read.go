package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/repository"
	"github.com/google/uuid"
)

func (repo *PostgresClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Client, error) {
	query := `
		SELECT id, full_name, rfc, phone, email, created_at
		FROM clients
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	var client entities.Client
	err := row.Scan(&client.ID, &client.FullName, &client.RFC, &client.Phone, &client.Email, &client.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (repo *PostgresClientRepository) GetByRFC(ctx context.Context, rfc string) (*entities.Client, error) {
	query := `
		SELECT id, full_name, rfc, phone, email, created_at
		FROM clients
		WHERE rfc = $1
	`
	row := repo.db.QueryRowContext(ctx, query, rfc)
	var client entities.Client
	err := row.Scan(&client.ID, &client.FullName, &client.RFC, &client.Phone, &client.Email, &client.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (repo *PostgresClientRepository) List(ctx context.Context, filters repository.ClientFilters) ([]*entities.Client, error) {
	baseQuery := `
		SELECT id, full_name, rfc, phone, email, created_at
		FROM clients
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.Search != nil && *filters.Search != "" {
		baseQuery += ` AND (full_name ILIKE $` + strconv.Itoa(argID) +
			` OR rfc ILIKE $` + strconv.Itoa(argID) +
			` OR email ILIKE $` + strconv.Itoa(argID) + `)`
		args = append(args, "%"+*filters.Search+"%")
		argID++
	}

	baseQuery += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*entities.Client
	for rows.Next() {
		var client entities.Client
		err := rows.Scan(&client.ID, &client.FullName, &client.RFC, &client.Phone, &client.Email, &client.CreatedAt)
		if err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}
	return clients, nil
}
