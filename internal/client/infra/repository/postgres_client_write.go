package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/entities"
	"github.com/google/uuid"
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

func (repo *PostgresClientRepository) CountWorksWithClientInStatus(ctx context.Context, clientID uuid.UUID, status string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM works WHERE client_id = $1 AND status = $2`
	err := repo.db.QueryRowContext(ctx, query, clientID, status).Scan(&count)
	return count, err
}

func (repo *PostgresClientRepository) UpdatePendingWorksClientID(ctx context.Context, oldClientID uuid.UUID, newClientID uuid.UUID) error {
	query := `UPDATE works SET client_id = $1 WHERE client_id = $2 AND status != 'APPROVED'`
	_, err := repo.db.ExecContext(ctx, query, newClientID, oldClientID)
	return err
}
