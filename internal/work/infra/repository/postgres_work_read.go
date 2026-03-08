package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

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
