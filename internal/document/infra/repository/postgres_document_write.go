package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/entities"
)

func (repo *PostgresDocumentRepository) Create(ctx context.Context, doc *entities.Document) error {
	query := `
		INSERT INTO documents (id, client_id, work_id, user_id, document_name, category, version, file_path, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := repo.db.ExecContext(ctx, query,
		doc.ID, doc.ClientID, doc.WorkID, doc.UserID,
		doc.DocumentName, doc.Category, doc.Version, doc.FilePath, doc.CreatedAt,
	)
	return err
}
