package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/entities"
)

func (repo *PostgresClientRepository) Create(ctx context.Context, client *entities.Client) error {
	query := `
		INSERT INTO clients (id, full_name, rfc, phone, email, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return repo.db.QueryRowContext(ctx, query,
		client.ID, client.FullName, client.RFC, client.Phone, client.Email, client.CreatedAt,
	).Scan(&client.ID)
}

func (repo *PostgresClientRepository) Update(ctx context.Context, client *entities.Client) error {
	query := `
		UPDATE clients
		SET full_name = $1, rfc = $2, phone = $3, email = $4
		WHERE id = $5
	`
	_, err := repo.db.ExecContext(ctx, query,
		client.FullName, client.RFC, client.Phone, client.Email, client.ID,
	)
	return err
}
