package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/entities"
	domainRepo "github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/repository"
	"github.com/google/uuid"
)

type PostgresAuditRepository struct {
	db *sql.DB
}

func NewPostgresAuditRepository(db *sql.DB) *PostgresAuditRepository {
	return &PostgresAuditRepository{db: db}
}

// ─── Create ─────────────────────────────────────────────────────────────────

func (r *PostgresAuditRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, action, entity, entity_id, json_details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.UserID,
		log.Action,
		log.Entity,
		log.EntityID,
		log.JSONDetails,
		log.CreatedAt,
	)
	return err
}

// ─── List (query dinámico) ──────────────────────────────────────────────────

func (r *PostgresAuditRepository) List(ctx context.Context, filters domainRepo.AuditFilters) ([]*entities.AuditLog, error) {
	baseQuery := `
		SELECT id, user_id, action, entity, entity_id, json_details, created_at
		FROM audit_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.UserID != nil && *filters.UserID != "" {
		baseQuery += ` AND user_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.UserID)
		argID++
	}

	if filters.Action != nil && *filters.Action != "" {
		baseQuery += ` AND action = $` + strconv.Itoa(argID)
		args = append(args, *filters.Action)
		argID++
	}

	if filters.Entity != nil && *filters.Entity != "" {
		baseQuery += ` AND entity = $` + strconv.Itoa(argID)
		args = append(args, *filters.Entity)
		argID++
	}

	if filters.EntityID != nil && *filters.EntityID != "" {
		baseQuery += ` AND entity_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.EntityID)
		argID++
	}

	if filters.StartDate != nil && *filters.StartDate != "" {
		baseQuery += ` AND created_at >= $` + strconv.Itoa(argID)
		args = append(args, *filters.StartDate)
		argID++
	}

	if filters.EndDate != nil && *filters.EndDate != "" {
		baseQuery += ` AND created_at <= $` + strconv.Itoa(argID)
		args = append(args, *filters.EndDate)
		argID++
	}

	baseQuery += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entities.AuditLog
	for rows.Next() {
		var l entities.AuditLog
		var userID *uuid.UUID

		if err := rows.Scan(
			&l.ID,
			&userID,
			&l.Action,
			&l.Entity,
			&l.EntityID,
			&l.JSONDetails,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		l.UserID = userID
		logs = append(logs, &l)
	}

	return logs, rows.Err()
}
