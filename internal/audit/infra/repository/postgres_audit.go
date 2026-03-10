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
		SELECT al.id, al.user_id, u.full_name as user_name, al.action, al.entity, al.entity_id, al.json_details, al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.UserID != nil && *filters.UserID != "" {
		baseQuery += ` AND al.user_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.UserID)
		argID++
	}

	if filters.Action != nil && *filters.Action != "" {
		baseQuery += ` AND al.action = $` + strconv.Itoa(argID)
		args = append(args, *filters.Action)
		argID++
	}

	if filters.Entity != nil && *filters.Entity != "" {
		baseQuery += ` AND al.entity = $` + strconv.Itoa(argID)
		args = append(args, *filters.Entity)
		argID++
	}

	if filters.EntityID != nil && *filters.EntityID != "" {
		baseQuery += ` AND al.entity_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.EntityID)
		argID++
	}

	if filters.StartDate != nil && *filters.StartDate != "" {
		baseQuery += ` AND al.created_at >= $` + strconv.Itoa(argID)
		args = append(args, *filters.StartDate)
		argID++
	}

	if filters.EndDate != nil && *filters.EndDate != "" {
		baseQuery += ` AND al.created_at <= $` + strconv.Itoa(argID)
		args = append(args, *filters.EndDate)
		argID++
	}

	baseQuery += ` ORDER BY al.created_at DESC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
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
		var userName *string
		var jsonDetails []byte

		if err := rows.Scan(
			&l.ID,
			&userID,
			&userName,
			&l.Action,
			&l.Entity,
			&l.EntityID,
			&jsonDetails,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		l.UserID = userID
		l.UserName = userName
		if len(jsonDetails) > 0 {
			l.JSONDetails = jsonDetails
		}
		logs = append(logs, &l)
	}

	return logs, rows.Err()
}

// ─── Count (para paginación real) ───────────────────────────────────────────

func (r *PostgresAuditRepository) Count(ctx context.Context, filters domainRepo.AuditFilters) (int, error) {
	baseQuery := `SELECT COUNT(*) FROM audit_logs al WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if filters.UserID != nil && *filters.UserID != "" {
		baseQuery += ` AND al.user_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.UserID)
		argID++
	}

	if filters.Action != nil && *filters.Action != "" {
		baseQuery += ` AND al.action = $` + strconv.Itoa(argID)
		args = append(args, *filters.Action)
		argID++
	}

	if filters.Entity != nil && *filters.Entity != "" {
		baseQuery += ` AND al.entity = $` + strconv.Itoa(argID)
		args = append(args, *filters.Entity)
		argID++
	}

	if filters.EntityID != nil && *filters.EntityID != "" {
		baseQuery += ` AND al.entity_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.EntityID)
		argID++
	}

	if filters.StartDate != nil && *filters.StartDate != "" {
		baseQuery += ` AND al.created_at >= $` + strconv.Itoa(argID)
		args = append(args, *filters.StartDate)
		argID++
	}

	if filters.EndDate != nil && *filters.EndDate != "" {
		baseQuery += ` AND al.created_at <= $` + strconv.Itoa(argID)
		args = append(args, *filters.EndDate)
	}

	var total int
	err := r.db.QueryRowContext(ctx, baseQuery, args...).Scan(&total)
	return total, err
}

// ─── Metrics ────────────────────────────────────────────────────────────────

func (r *PostgresAuditRepository) GetUserActionMetrics(ctx context.Context, filters domainRepo.AuditFilters) ([]entities.ActionMetricItem, error) {
	baseQuery := `
		SELECT al.action, COUNT(al.id) as count
		FROM audit_logs al
		WHERE al.entity = 'USER'
	`
	args := []interface{}{}
	argID := 1

	if filters.StartDate != nil && *filters.StartDate != "" {
		baseQuery += ` AND al.created_at >= $` + strconv.Itoa(argID)
		args = append(args, *filters.StartDate)
		argID++
	}

	if filters.EndDate != nil && *filters.EndDate != "" {
		baseQuery += ` AND al.created_at <= $` + strconv.Itoa(argID)
		args = append(args, *filters.EndDate)
		argID++
	}

	baseQuery += ` GROUP BY al.action ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []entities.ActionMetricItem
	for rows.Next() {
		var item entities.ActionMetricItem
		if err := rows.Scan(&item.Action, &item.Count); err != nil {
			return nil, err
		}
		metrics = append(metrics, item)
	}

	return metrics, rows.Err()
}

func (r *PostgresAuditRepository) GetWorkActionMetrics(ctx context.Context, filters domainRepo.AuditFilters) ([]entities.ActionMetricItem, error) {
	// Group actions where entity is WORK
	baseQuery := `
		SELECT 
			CASE 
				WHEN action = 'STATUS_CHANGE' THEN COALESCE(json_details->>'new_status', 'STATUS_CHANGE')
				ELSE action 
			END as mapped_action, 
			COUNT(id) as count
		FROM audit_logs
		WHERE entity = 'WORK'
	`
	args := []interface{}{}
	argID := 1

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

	baseQuery += ` GROUP BY mapped_action ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []entities.ActionMetricItem
	for rows.Next() {
		var item entities.ActionMetricItem
		if err := rows.Scan(&item.Action, &item.Count); err != nil {
			return nil, err
		}
		metrics = append(metrics, item)
	}

	return metrics, rows.Err()
}
