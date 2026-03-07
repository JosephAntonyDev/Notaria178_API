package repository

import (
	"context"
	"database/sql"
	"strconv"

	domainRepo "github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/domain/repository"
)

type PostgresDashboardRepository struct {
	db *sql.DB
}

func NewPostgresDashboardRepository(db *sql.DB) *PostgresDashboardRepository {
	return &PostgresDashboardRepository{db: db}
}

// ─── KPIs: conteos por estado en una sola query ─────────────────────────────

func (r *PostgresDashboardRepository) GetKPIs(ctx context.Context, filters domainRepo.DashboardFilters) (*domainRepo.KPIsResult, error) {
	query := `
		SELECT
			COUNT(*)                                            AS total,
			COUNT(*) FILTER (WHERE status = 'PENDING')          AS pending,
			COUNT(*) FILTER (WHERE status = 'IN_PROGRESS')      AS in_progress,
			COUNT(*) FILTER (WHERE status = 'READY_FOR_REVIEW') AS ready_for_review,
			COUNT(*) FILTER (WHERE status = 'APPROVED')         AS approved,
			COUNT(*) FILTER (WHERE status = 'REJECTED')         AS rejected
		FROM works
		WHERE created_at >= $1 AND created_at < $2
	`
	args := []interface{}{filters.StartDate, filters.EndDate}
	argID := 3

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		query += ` AND branch_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.BranchID)
	}

	var result domainRepo.KPIsResult
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&result.Total,
		&result.Pending,
		&result.InProgress,
		&result.ReadyForReview,
		&result.Approved,
		&result.Rejected,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ─── Trend: trabajos ingresados vs aprobados agrupados por periodo ──────────

func (r *PostgresDashboardRepository) GetTrend(ctx context.Context, filters domainRepo.DashboardFilters, groupBy string) ([]domainRepo.TrendRow, error) {
	// PostgreSQL date_trunc acepta 'day', 'week', 'month' directamente.
	// Validamos en el use case; aquí solo lo usamos.
	query := `
		SELECT
			date_trunc('` + groupBy + `', created_at)::date AS period,
			COUNT(*)                                        AS created,
			COUNT(*) FILTER (WHERE status = 'APPROVED')     AS approved
		FROM works
		WHERE created_at >= $1 AND created_at < $2
	`
	args := []interface{}{filters.StartDate, filters.EndDate}
	argID := 3

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		query += ` AND branch_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.BranchID)
	}

	query += ` GROUP BY period ORDER BY period ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domainRepo.TrendRow
	for rows.Next() {
		var row domainRepo.TrendRow
		if err := rows.Scan(&row.Period, &row.Created, &row.Approved); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ─── Distribution: trabajos agrupados por estado ────────────────────────────

func (r *PostgresDashboardRepository) GetDistribution(ctx context.Context, filters domainRepo.DashboardFilters) ([]domainRepo.DistributionRow, error) {
	query := `
		SELECT status::text, COUNT(*) AS count
		FROM works
		WHERE created_at >= $1 AND created_at < $2
	`
	args := []interface{}{filters.StartDate, filters.EndDate}
	argID := 3

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		query += ` AND branch_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.BranchID)
	}

	query += ` GROUP BY status ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domainRepo.DistributionRow
	for rows.Next() {
		var row domainRepo.DistributionRow
		if err := rows.Scan(&row.Status, &row.Count); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ─── Activity: actividad reciente con nombre de usuario ─────────────────────

func (r *PostgresDashboardRepository) GetRecentActivity(ctx context.Context, filters domainRepo.ActivityFilters) ([]domainRepo.ActivityRow, int, error) {
	// ── Count total (para paginación) ───────────────────────────────────
	countQuery := `
		SELECT COUNT(*)
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.created_at >= $1 AND al.created_at < $2
	`
	args := []interface{}{filters.StartDate, filters.EndDate}
	argID := 3

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		countQuery += ` AND u.branch_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.BranchID)
		argID++
	}
	if filters.UserID != nil && *filters.UserID != "" {
		countQuery += ` AND al.user_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.UserID)
		argID++
	}
	if filters.EntityID != nil && *filters.EntityID != "" {
		countQuery += ` AND al.entity_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.EntityID)
		argID++
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// ── Data query ──────────────────────────────────────────────────────
	dataQuery := `
		SELECT al.id, al.user_id, u.full_name,
		       al.action, al.entity, al.entity_id, al.json_details, al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.created_at >= $1 AND al.created_at < $2
	`
	dataArgs := []interface{}{filters.StartDate, filters.EndDate}
	dataArgID := 3

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		dataQuery += ` AND u.branch_id = $` + strconv.Itoa(dataArgID)
		dataArgs = append(dataArgs, *filters.BranchID)
		dataArgID++
	}
	if filters.UserID != nil && *filters.UserID != "" {
		dataQuery += ` AND al.user_id = $` + strconv.Itoa(dataArgID)
		dataArgs = append(dataArgs, *filters.UserID)
		dataArgID++
	}
	if filters.EntityID != nil && *filters.EntityID != "" {
		dataQuery += ` AND al.entity_id = $` + strconv.Itoa(dataArgID)
		dataArgs = append(dataArgs, *filters.EntityID)
		dataArgID++
	}

	dataQuery += ` ORDER BY al.created_at DESC LIMIT $` + strconv.Itoa(dataArgID) + ` OFFSET $` + strconv.Itoa(dataArgID+1)
	dataArgs = append(dataArgs, filters.Limit, filters.Offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []domainRepo.ActivityRow
	for rows.Next() {
		var row domainRepo.ActivityRow
		if err := rows.Scan(
			&row.ID, &row.UserID, &row.UserName,
			&row.Action, &row.Entity, &row.EntityID, &row.JSONDetails, &row.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		results = append(results, row)
	}
	return results, total, rows.Err()
}

// ─── Top Drafters: proyectistas con más trabajos asignados ──────────────────

func (r *PostgresDashboardRepository) GetTopDrafters(ctx context.Context, filters domainRepo.DashboardFilters, limit int) ([]domainRepo.TopDrafterRow, error) {
	query := `
		SELECT u.id::text, u.full_name, u.role::text, COUNT(w.id) AS work_count
		FROM users u
		JOIN works w ON w.main_drafter_id = u.id
		WHERE w.created_at >= $1 AND w.created_at < $2
	`
	args := []interface{}{filters.StartDate, filters.EndDate}
	argID := 3

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		query += ` AND u.branch_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.BranchID)
		argID++
	}

	query += ` GROUP BY u.id, u.full_name, u.role ORDER BY work_count DESC LIMIT $` + strconv.Itoa(argID)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domainRepo.TopDrafterRow
	for rows.Next() {
		var row domainRepo.TopDrafterRow
		if err := rows.Scan(&row.UserID, &row.FullName, &row.Role, &row.WorkCount); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// ─── Top Acts: actos más frecuentes en expedientes ──────────────────────────

func (r *PostgresDashboardRepository) GetTopActs(ctx context.Context, filters domainRepo.DashboardFilters, limit int) ([]domainRepo.TopActRow, error) {
	query := `
		SELECT ac.id::text, ac.name, COUNT(wa.work_id) AS count
		FROM work_acts wa
		JOIN act_catalogs ac ON wa.act_id = ac.id
		JOIN works w ON wa.work_id = w.id
		WHERE w.created_at >= $1 AND w.created_at < $2
	`
	args := []interface{}{filters.StartDate, filters.EndDate}
	argID := 3

	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		query += ` AND w.branch_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.BranchID)
		argID++
	}

	query += ` GROUP BY ac.id, ac.name ORDER BY count DESC LIMIT $` + strconv.Itoa(argID)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domainRepo.TopActRow
	for rows.Next() {
		var row domainRepo.TopActRow
		if err := rows.Scan(&row.ActID, &row.Name, &row.Count); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}
