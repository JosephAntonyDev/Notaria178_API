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
		SELECT id, client_id, work_id, user_id, document_name, category, version, file_path, created_at
		FROM documents
		WHERE id = $1
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	var doc entities.Document
	err := row.Scan(
		&doc.ID, &doc.ClientID, &doc.WorkID, &doc.UserID,
		&doc.DocumentName, &doc.Category, &doc.Version, &doc.FilePath, &doc.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (repo *PostgresDocumentRepository) GetByWorkID(ctx context.Context, workID uuid.UUID) ([]*entities.Document, error) {
	query := `
		SELECT id, client_id, work_id, user_id, document_name, category, version, file_path, created_at
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
		if err := rows.Scan(
			&doc.ID, &doc.ClientID, &doc.WorkID, &doc.UserID,
			&doc.DocumentName, &doc.Category, &doc.Version, &doc.FilePath, &doc.CreatedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}
	return docs, nil
}
