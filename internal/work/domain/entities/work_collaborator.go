package entities

import "github.com/google/uuid"

// ─── Value Object para JOIN work_collaborators ↔ users (tabla: work_collaborators) ──

type WorkCollaboratorInfo struct {
	UserID   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
}
