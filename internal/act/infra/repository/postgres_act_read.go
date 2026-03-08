package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/google/uuid"
)

func (repo *PostgresActRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Act, error) {
	query := `
		SELECT a.id, a.name, a.description, a.category, a.status,
		       COUNT(DISTINCT r.id) as requirements_count,
		       COUNT(DISTINCT w.work_id) as works_count
		FROM act_catalogs a
		LEFT JOIN act_requirements r ON a.id = r.act_id AND r.status = 'ACTIVE'
		LEFT JOIN work_acts w ON a.id = w.act_id
		WHERE a.id = $1
		GROUP BY a.id, a.name, a.description, a.category, a.status
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	var act entities.Act
	var category sql.NullString
	var desc sql.NullString
	err := row.Scan(
		&act.ID, &act.Name, &desc, &category, &act.Status,
		&act.RequirementsCount, &act.WorksCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("[ERROR] GetByID Scan failed: %v", err)
		return nil, err
	}
	if desc.Valid {
		act.Description = &desc.String
	}
	if category.Valid {
		act.Category = category.String
	} else {
		act.Category = "General"
	}
	return &act, nil
}

func (repo *PostgresActRepository) GetByName(ctx context.Context, name string) (*entities.Act, error) {
	query := `
		SELECT a.id, a.name, a.description, a.category, a.status,
		       COUNT(DISTINCT r.id) as requirements_count,
		       COUNT(DISTINCT w.work_id) as works_count
		FROM act_catalogs a
		LEFT JOIN act_requirements r ON a.id = r.act_id AND r.status = 'ACTIVE'
		LEFT JOIN work_acts w ON a.id = w.act_id
		WHERE a.name = $1
		GROUP BY a.id, a.name, a.description, a.category, a.status
	`
	row := repo.db.QueryRowContext(ctx, query, name)
	var act entities.Act
	var category sql.NullString
	var desc sql.NullString
	err := row.Scan(
		&act.ID, &act.Name, &desc, &category, &act.Status,
		&act.RequirementsCount, &act.WorksCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("[ERROR] GetByName Scan failed: %v", err)
		return nil, err
	}
	if desc.Valid {
		act.Description = &desc.String
	}
	if category.Valid {
		act.Category = category.String
	} else {
		act.Category = "General"
	}
	return &act, nil
}

func (repo *PostgresActRepository) List(ctx context.Context, filters repository.ActFilters) ([]*entities.Act, error) {
	baseQuery := `
		SELECT a.id, a.name, a.description, a.category, a.status,
		       COUNT(DISTINCT r.id) as requirements_count,
		       COUNT(DISTINCT w.work_id) as works_count
		FROM act_catalogs a
		LEFT JOIN act_requirements r ON a.id = r.act_id AND r.status = 'ACTIVE'
		LEFT JOIN work_acts w ON a.id = w.act_id
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if filters.Search != nil && *filters.Search != "" {
		baseQuery += ` AND a.name ILIKE $` + strconv.Itoa(argID)
		args = append(args, "%"+*filters.Search+"%")
		argID++
	}
	if filters.Status != nil && *filters.Status != "" {
		baseQuery += ` AND a.status = $` + strconv.Itoa(argID)
		args = append(args, *filters.Status)
		argID++
	}

	baseQuery += ` GROUP BY a.id, a.name, a.description, a.category, a.status ORDER BY a.category ASC, a.name ASC LIMIT $` + strconv.Itoa(argID) + ` OFFSET $` + strconv.Itoa(argID+1)
	args = append(args, filters.Limit, filters.Offset)

	log.Printf("[DEBUG] List SQL: %s | args: %v", baseQuery, args)
	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		log.Printf("[ERROR] List QueryContext failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var acts []*entities.Act
	for rows.Next() {
		var act entities.Act
		var category sql.NullString
		var desc sql.NullString
		
		err := rows.Scan(
			&act.ID, &act.Name, &desc, &category, &act.Status,
			&act.RequirementsCount, &act.WorksCount,
		)
		if err != nil {
			log.Printf("[ERROR] List Scan failed: %v", err)
			return nil, err
		}
		if desc.Valid {
			act.Description = &desc.String
		}
		if category.Valid {
			act.Category = category.String
		} else {
			act.Category = "General"
		}
		acts = append(acts, &act)
	}
	return acts, nil
}
