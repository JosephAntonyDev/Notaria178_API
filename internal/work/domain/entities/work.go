package entities

import (
	"time"

	"github.com/google/uuid"
)

// ─── Status del expediente ───────────────────────────────────────────────────

type WorkStatus string

const (
	StatusPending        WorkStatus = "PENDING"
	StatusInProgress     WorkStatus = "IN_PROGRESS"
	StatusReadyForReview WorkStatus = "READY_FOR_REVIEW"
	StatusRejected       WorkStatus = "REJECTED"
	StatusApproved       WorkStatus = "APPROVED"
)

// ─── Entidad principal (tabla: works) ───────────────────────────────────────

type Work struct {
	ID            uuid.UUID  `json:"id"`
	BranchID      uuid.UUID  `json:"branch_id"`
	ClientID      uuid.UUID  `json:"client_id"`
	MainDrafterID *uuid.UUID `json:"main_drafter_id,omitempty"`
	Folio         *string    `json:"folio,omitempty"`
	Status        WorkStatus `json:"status"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
