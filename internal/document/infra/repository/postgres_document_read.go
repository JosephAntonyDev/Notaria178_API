package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Document, error) {
	query := `
		SELECT id, client_id, work_id, user_id, document_name, category, version, file_path, requirement_id, requirement_source, created_at
		FROM documents
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	var doc entities.Document
	var reqSource sql.NullString
	var nullClientID, nullWorkID, nullUserID, nullReqID uuid.NullUUID

	err := row.Scan(
		&doc.ID, &nullClientID, &nullWorkID, &nullUserID,
		&doc.DocumentName, &doc.Category, &doc.Version, &doc.FilePath,
		&nullReqID, &reqSource, &doc.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if nullClientID.Valid {
		doc.ClientID = &nullClientID.UUID
	}
	if nullWorkID.Valid {
		doc.WorkID = &nullWorkID.UUID
	}
	if nullUserID.Valid {
		doc.UserID = &nullUserID.UUID
	}
	if nullReqID.Valid {
		doc.RequirementID = &nullReqID.UUID
	}
	if reqSource.Valid {
		doc.RequirementSource = reqSource.String
	}
	return &doc, nil
}

func (repo *PostgresDocumentRepository) GetByWorkID(ctx context.Context, workID uuid.UUID) ([]*entities.Document, error) {
	query := `
		SELECT id, client_id, work_id, user_id, document_name, category, version, file_path, requirement_id, requirement_source, created_at
		FROM documents
		WHERE work_id = $1
		ORDER BY created_at DESC
	`
	rows, err := repo.db.QueryContext(ctx, query, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := make([]*entities.Document, 0)
	for rows.Next() {
		var doc entities.Document
		var reqSource sql.NullString
		var nullClientID, nullWorkID, nullUserID, nullReqID uuid.NullUUID

		if err := rows.Scan(
			&doc.ID, &nullClientID, &nullWorkID, &nullUserID,
			&doc.DocumentName, &doc.Category, &doc.Version, &doc.FilePath,
			&nullReqID, &reqSource, &doc.CreatedAt,
		); err != nil {
			return nil, err
		}

		if nullClientID.Valid {
			doc.ClientID = &nullClientID.UUID
		}
		if nullWorkID.Valid {
			doc.WorkID = &nullWorkID.UUID
		}
		if nullUserID.Valid {
			doc.UserID = &nullUserID.UUID
		}
		if nullReqID.Valid {
			doc.RequirementID = &nullReqID.UUID
		}
		if reqSource.Valid {
			doc.RequirementSource = reqSource.String
		}
		docs = append(docs, &doc)
	}
	return docs, nil
}
