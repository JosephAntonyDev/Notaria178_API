package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"

	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func (repo *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT 
			id, branch_id, full_name, email, password_hash, 
			phone, role, status, hire_date, start_time, end_time, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	row := repo.db.QueryRowContext(ctx, query, email)

	var user entities.User

	err := row.Scan(
		&user.ID,
		&user.BranchID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Phone,
		&user.Role,
		&user.Status,
		&user.HireDate,
		&user.StartTime,
		&user.EndTime,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil 
		}
		return nil, err
	}

	return &user, nil
}

func (repo *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT 
			id, branch_id, full_name, email, password_hash, 
			phone, role, status, hire_date, start_time, end_time, created_at, updated_at
		FROM users 
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)

	var user entities.User
	err := row.Scan(
		&user.ID, &user.BranchID, &user.FullName, &user.Email,
		&user.PasswordHash, &user.Phone, &user.Role, &user.Status,
		&user.HireDate, &user.StartTime, &user.EndTime, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresUserRepository) List(ctx context.Context, limit int, offset int) ([]*entities.User, error) {
	query := `
		SELECT 
			id, branch_id, full_name, email, password_hash, 
			phone, role, status, hire_date, start_time, end_time, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := repo.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
			&user.ID, &user.BranchID, &user.FullName, &user.Email,
			&user.PasswordHash, &user.Phone, &user.Role, &user.Status,
			&user.HireDate, &user.StartTime, &user.EndTime, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}