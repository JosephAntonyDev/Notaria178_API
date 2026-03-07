package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"strconv"

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

func (repo *PostgresUserRepository) List(ctx context.Context, filters entities.UserFilters) ([]*entities.User, int, error) {
	
	baseQuery := `
		SELECT 
			id, branch_id, full_name, email, password_hash, 
			phone, role, status, hire_date, start_time, end_time, created_at, updated_at
		FROM users 
		WHERE 1=1 
	`
	
	args := []interface{}{}
	argId := 1

	if filters.Search != nil && *filters.Search != "" {
		baseQuery += ` AND (full_name ILIKE $` + strconv.Itoa(argId) + ` OR email ILIKE $` + strconv.Itoa(argId) + `)`
		args = append(args, "%"+*filters.Search+"%")
		argId++
	}

	if filters.Status != nil && *filters.Status != "" {
		baseQuery += ` AND status = $` + strconv.Itoa(argId)
		args = append(args, *filters.Status)
		argId++
	}

	if filters.Role != nil && *filters.Role != "" {
		baseQuery += ` AND role = $` + strconv.Itoa(argId)
		args = append(args, *filters.Role)
		argId++
	}

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		baseQuery += ` AND branch_id = $` + strconv.Itoa(argId)
		args = append(args, *filters.BranchID)
		argId++
	}

	if filters.StartDate != nil && *filters.StartDate != "" {
		baseQuery += ` AND hire_date >= $` + strconv.Itoa(argId)
		args = append(args, *filters.StartDate)
		argId++
	}

	if filters.EndDate != nil && *filters.EndDate != "" {
		baseQuery += ` AND hire_date <= $` + strconv.Itoa(argId)
		args = append(args, *filters.EndDate)
		argId++
	}

	orderDir := "DESC"
	if filters.Sort != nil && *filters.Sort == "asc" {
		orderDir = "ASC"
	}

	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") AS sub"
	var total int
	err := repo.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	baseQuery += ` ORDER BY created_at ` + orderDir + ` LIMIT $` + strconv.Itoa(argId) + ` OFFSET $` + strconv.Itoa(argId+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		users = append(users, &user)
	}
	return users, total, nil
}