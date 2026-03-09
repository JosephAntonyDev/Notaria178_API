package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	importLog "log"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

func (repo *PostgresWorkRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Work, error) {
	query := `
		SELECT id, branch_id, client_id, main_drafter_id, folio, status, deadline, created_at, updated_at
		FROM works
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	var work entities.Work
	var branchID uuid.NullUUID
	err := row.Scan(
		&work.ID, &branchID, &work.ClientID, &work.MainDrafterID,
		&work.Folio, &work.Status, &work.Deadline, &work.CreatedAt, &work.UpdatedAt,
	)
	if branchID.Valid {
		work.BranchID = branchID.UUID
	} else {
		work.BranchID = uuid.Nil
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &work, nil
}

func (repo *PostgresWorkRepository) List(ctx context.Context, filters repository.WorkFilters) ([]*entities.Work, error) {
	baseQuery := `
		SELECT id, branch_id, client_id, main_drafter_id, folio, status, deadline, created_at, updated_at
		FROM works
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.Search != nil && *filters.Search != "" {
		baseQuery += ` AND folio ILIKE $` + strconv.Itoa(argID)
		args = append(args, "%"+*filters.Search+"%")
		argID++
	}
	if filters.Status != nil && *filters.Status != "" {
		baseQuery += ` AND status = $` + strconv.Itoa(argID)
		args = append(args, *filters.Status)
		argID++
	}
	if filters.BranchID != nil && *filters.BranchID != "" && *filters.BranchID != "all" {
		baseQuery += ` AND branch_id = $` + strconv.Itoa(argID)
		args = append(args, *filters.BranchID)
		argID++
	}
	if filters.ScopedUserID != nil && *filters.ScopedUserID != "" {
		baseQuery += ` AND (main_drafter_id = $` + strconv.Itoa(argID) +
			` OR id IN (SELECT work_id FROM work_collaborators WHERE user_id = $` + strconv.Itoa(argID+1) + `))`
		args = append(args, *filters.ScopedUserID, *filters.ScopedUserID)
		argID += 2
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

	orderDir := "DESC "
	if filters.Sort != nil && *filters.Sort == "asc" {
		orderDir = "ASC "
	}

	baseQuery += ` ORDER BY updated_at ` + orderDir + `LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		importLog.Printf("SQL Error in Works List QueryContext: %v\nQuery: %s\nArgs: %v", err, baseQuery, args)
		return nil, err
	}
	defer rows.Close()

	var works []*entities.Work
	for rows.Next() {
		var work entities.Work
		var branchID uuid.NullUUID
		err := rows.Scan(
			&work.ID, &branchID, &work.ClientID, &work.MainDrafterID,
			&work.Folio, &work.Status, &work.Deadline, &work.CreatedAt, &work.UpdatedAt,
		)
		if branchID.Valid {
			work.BranchID = branchID.UUID
		} else {
			work.BranchID = uuid.Nil
		}
		if err != nil {
			importLog.Printf("SQL Error in Works List Scan: %v", err)
			return nil, err
		}
		works = append(works, &work)
	}
	return works, nil
}

func (repo *PostgresWorkRepository) GetClientNameByID(ctx context.Context, clientID uuid.UUID) (string, error) {
	var name string
	err := repo.db.QueryRowContext(ctx, `SELECT full_name FROM clients WHERE id = $1`, clientID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (repo *PostgresWorkRepository) GetBranchNameByID(ctx context.Context, branchID uuid.UUID) (string, error) {
	var name string
	err := repo.db.QueryRowContext(ctx, `SELECT name FROM branches WHERE id = $1`, branchID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (repo *PostgresWorkRepository) GetUserFullNameByID(ctx context.Context, userID uuid.UUID) (string, error) {
	var name string
	err := repo.db.QueryRowContext(ctx, `SELECT full_name FROM users WHERE id = $1`, userID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (repo *PostgresWorkRepository) GetClientByID(ctx context.Context, clientID uuid.UUID) (*entities.ClientInfo, error) {
	query := `SELECT id, full_name, rfc, phone, email FROM clients WHERE id = $1`
	row := repo.db.QueryRowContext(ctx, query, clientID)
	var c entities.ClientInfo
	err := row.Scan(&c.ID, &c.FullName, &c.RFC, &c.Phone, &c.Email)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (repo *PostgresWorkRepository) CountWorksWithClientInStatus(ctx context.Context, clientID uuid.UUID, status string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM works WHERE client_id = $1 AND status = $2`
	err := repo.db.QueryRowContext(ctx, query, clientID, status).Scan(&count)
	return count, err
}

func (repo *PostgresWorkRepository) GetRequirementsByActIDs(ctx context.Context, actIDs []uuid.UUID) ([]entities.ActRequirementInfo, error) {
	if len(actIDs) == 0 {
		return nil, nil
	}
	// Build placeholders
	args := make([]interface{}, len(actIDs))
	placeholders := ""
	for i, id := range actIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "$" + strconv.Itoa(i+1)
		args[i] = id
	}
	query := `
		SELECT id, act_id, name, status
		FROM act_requirements
		WHERE act_id IN (` + placeholders + `) AND status = 'ACTIVE'
		ORDER BY name ASC
	`
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reqs := make([]entities.ActRequirementInfo, 0)
	for rows.Next() {
		var r entities.ActRequirementInfo
		if err := rows.Scan(&r.ID, &r.ActID, &r.Name, &r.Status); err != nil {
			return nil, err
		}
		reqs = append(reqs, r)
	}
	return reqs, nil
}

func (repo *PostgresWorkRepository) GetRequirementDocumentsByWorkID(ctx context.Context, workID uuid.UUID) (map[uuid.UUID]uuid.UUID, error) {
	query := `
		SELECT DISTINCT ON (requirement_id) requirement_id, id
		FROM documents
		WHERE work_id = $1 AND category = 'CLIENT_REQUIREMENT' AND requirement_id IS NOT NULL
		ORDER BY requirement_id, created_at DESC
	`
	rows, err := repo.db.QueryContext(ctx, query, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]uuid.UUID)
	for rows.Next() {
		var reqID, docID uuid.UUID
		if err := rows.Scan(&reqID, &docID); err != nil {
			return nil, err
		}
		result[reqID] = docID
	}
	return result, nil
}

func (repo *PostgresWorkRepository) GetDocumentsForCleanupByReqIDs(ctx context.Context, workID uuid.UUID, reqIDs []uuid.UUID) ([]entities.DocCleanupInfo, error) {
	if len(reqIDs) == 0 {
		return nil, nil
	}
	placeholders := ""
	args := []interface{}{workID}
	for i, id := range reqIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "$" + strconv.Itoa(i+2)
		args = append(args, id)
	}
	query := `
		SELECT id, file_path, document_name
		FROM documents
		WHERE work_id = $1 AND category = 'CLIENT_REQUIREMENT' AND requirement_id IN (` + placeholders + `)
	`
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := make([]entities.DocCleanupInfo, 0)
	for rows.Next() {
		var d entities.DocCleanupInfo
		if err := rows.Scan(&d.ID, &d.FilePath, &d.Name); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

func (repo *PostgresWorkRepository) GetActRequirementIDsByNames(ctx context.Context, names []string) ([]uuid.UUID, error) {
	if len(names) == 0 {
		return nil, nil
	}
	placeholders := ""
	args := make([]interface{}, len(names))
	for i, n := range names {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "$" + strconv.Itoa(i+1)
		args[i] = strings.ToLower(strings.TrimSpace(n))
	}
	query := `SELECT id FROM act_requirements WHERE LOWER(TRIM(name)) IN (` + placeholders + `)`
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
