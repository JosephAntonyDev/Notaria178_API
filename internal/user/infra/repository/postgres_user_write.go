package repository

import (
	"context"
	"time"
	"github.com/google/uuid"


	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func (repo *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (
			id, branch_id, full_name, email, password_hash, 
			phone, role, status, hire_date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	_, err := repo.db.ExecContext(ctx, query,
		user.ID,
		user.BranchID,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.Role,
		user.Status,
		user.HireDate,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET full_name = $1, phone = $2, role = $3, branch_id = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := repo.db.ExecContext(ctx, query,
		user.FullName, user.Phone, user.Role, user.BranchID, time.Now(), user.ID,
	)
	return err
}

func (repo *PostgresUserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entities.UserStatus) error {
	query := `
		UPDATE users 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := repo.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}