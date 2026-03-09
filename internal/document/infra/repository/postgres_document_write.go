package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/entities"
	"github.com/google/uuid"
)

func (repo *PostgresDocumentRepository) Create(ctx context.Context, doc *entities.Document) error {
	query := `
		INSERT INTO documents (id, client_id, work_id, user_id, document_name, category, version, file_path, requirement_id, requirement_source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := repo.db.ExecContext(ctx, query,
		doc.ID, doc.ClientID, doc.WorkID, doc.UserID,
		doc.DocumentName, doc.Category, doc.Version, doc.FilePath,
		doc.RequirementID, doc.RequirementSource, doc.CreatedAt,
	)
	return err
}

func (repo *PostgresDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *PostgresDocumentRepository) LinkDocumentToWorkRequirement(ctx context.Context, docID uuid.UUID, reqID uuid.UUID) error {
	query := `UPDATE work_requirements SET document_id = $1 WHERE id = $2`
	_, err := repo.db.ExecContext(ctx, query, docID, reqID)
	return err
}
