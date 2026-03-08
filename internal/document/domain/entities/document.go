package entities

import (
	"time"

	"github.com/google/uuid"
)

// ─── Enum de categoría de documento ─────────────────────────────────────────

type DocumentCategory string

const (
	CategoryDraftDeed         DocumentCategory = "DRAFT_DEED"
	CategoryFinalDeed         DocumentCategory = "FINAL_DEED"
	CategoryClientRequirement DocumentCategory = "CLIENT_REQUIREMENT"
	CategoryOther             DocumentCategory = "OTHER"
)

// ─── Entidad principal (tabla: documents) ───────────────────────────────────

type Document struct {
	ID           uuid.UUID        `json:"id"`
	ClientID     *uuid.UUID       `json:"client_id,omitempty"`
	WorkID       *uuid.UUID       `json:"work_id,omitempty"`
	UserID       *uuid.UUID       `json:"user_id,omitempty"`
	DocumentName string           `json:"document_name"`
	Category     DocumentCategory `json:"category"`
	Version      int              `json:"version"`
	FilePath     string           `json:"file_path"`
	CreatedAt    time.Time        `json:"created_at"`
}
