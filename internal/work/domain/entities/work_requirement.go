package entities

import (
	"time"

	"github.com/google/uuid"
)

// ClientInfo contiene datos completos del cliente para enriquecer la vista de trabajo
type ClientInfo struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	RFC      *string   `json:"rfc,omitempty"`
	Phone    *string   `json:"phone,omitempty"`
	Email    *string   `json:"email,omitempty"`
}

// ActRequirementInfo contiene información de un requisito de acto con su act_id
type ActRequirementInfo struct {
	ID     uuid.UUID `json:"id"`
	ActID  uuid.UUID `json:"act_id"`
	Name   string    `json:"name"`
	Status string    `json:"status"`
}

// WorkRequirement es un requisito extra/ad-hoc asociado directamente al trabajo
type WorkRequirement struct {
	ID         uuid.UUID  `json:"id"`
	WorkID     uuid.UUID  `json:"work_id"`
	Name       string     `json:"name"`
	DocumentID *uuid.UUID `json:"document_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// DocCleanupInfo contiene la info mínima para borrar un documento del disco y BD
type DocCleanupInfo struct {
	ID       uuid.UUID
	FilePath string
	Name     string
}
