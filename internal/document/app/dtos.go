package app

import (
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/domain/entities"
	"github.com/google/uuid"
)

// ─── DTOs ───────────────────────────────────────────────────────────────────

type DocumentDTO struct {
	ID           uuid.UUID  `json:"id"`
	ClientID     *uuid.UUID `json:"client_id,omitempty"`
	WorkID       *uuid.UUID `json:"work_id,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	DocumentName string     `json:"document_name"`
	Category     string     `json:"category"`
	Version      int        `json:"version"`
	FilePath     string     `json:"file_path"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ─── Mapper ─────────────────────────────────────────────────────────────────

func ToDocumentDTO(doc *entities.Document) DocumentDTO {
	return DocumentDTO{
		ID:           doc.ID,
		ClientID:     doc.ClientID,
		WorkID:       doc.WorkID,
		UserID:       doc.UserID,
		DocumentName: doc.DocumentName,
		Category:     string(doc.Category),
		Version:      doc.Version,
		FilePath:     doc.FilePath,
		CreatedAt:    doc.CreatedAt,
	}
}
