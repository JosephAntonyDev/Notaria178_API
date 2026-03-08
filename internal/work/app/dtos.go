package app

import (
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

// ─── Contexto del request (extraído del JWT) ────────────────────────────────

type RequestContext struct {
	UserID   string
	UserRole string
	BranchID string // vacío para SUPER_ADMIN
}

// ─── DTOs ───────────────────────────────────────────────────────────────────

type WorkDTO struct {
	ID            uuid.UUID  `json:"id"`
	BranchID      uuid.UUID  `json:"branch_id"`
	ClientID      uuid.UUID  `json:"client_id"`
	MainDrafterID *uuid.UUID `json:"main_drafter_id,omitempty"`
	Folio         *string    `json:"folio,omitempty"`
	Status        string     `json:"status"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type WorkDetailDTO struct {
	WorkDTO
	Acts          []WorkActInfoDTO      `json:"acts"`
	Collaborators []WorkCollaboratorDTO `json:"collaborators"`
}

type WorkActInfoDTO struct {
	ActID uuid.UUID `json:"act_id"`
	Name  string    `json:"name"`
}

type WorkCollaboratorDTO struct {
	UserID   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
}

type WorkCommentDTO struct {
	ID        uuid.UUID `json:"id"`
	WorkID    uuid.UUID `json:"work_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// ─── Mappers ────────────────────────────────────────────────────────────────

func ToWorkDTO(work *entities.Work) WorkDTO {
	return WorkDTO{
		ID:            work.ID,
		BranchID:      work.BranchID,
		ClientID:      work.ClientID,
		MainDrafterID: work.MainDrafterID,
		Folio:         work.Folio,
		Status:        string(work.Status),
		Deadline:      work.Deadline,
		CreatedAt:     work.CreatedAt,
		UpdatedAt:     work.UpdatedAt,
	}
}

func ToWorkActInfoDTO(info entities.WorkActInfo) WorkActInfoDTO {
	return WorkActInfoDTO{ActID: info.ActID, Name: info.Name}
}

func ToWorkCollaboratorDTO(info entities.WorkCollaboratorInfo) WorkCollaboratorDTO {
	return WorkCollaboratorDTO{UserID: info.UserID, FullName: info.FullName}
}

func ToWorkCommentDTO(c entities.WorkComment) WorkCommentDTO {
	return WorkCommentDTO{
		ID: c.ID, WorkID: c.WorkID, UserID: c.UserID,
		UserName: c.UserName, Message: c.Message, CreatedAt: c.CreatedAt,
	}
}

func buildWorkDetail(work *entities.Work, acts []entities.WorkActInfo, collabs []entities.WorkCollaboratorInfo) *WorkDetailDTO {
	actsDTO := make([]WorkActInfoDTO, 0, len(acts))
	for _, a := range acts {
		actsDTO = append(actsDTO, ToWorkActInfoDTO(a))
	}
	collabsDTO := make([]WorkCollaboratorDTO, 0, len(collabs))
	for _, c := range collabs {
		collabsDTO = append(collabsDTO, ToWorkCollaboratorDTO(c))
	}
	return &WorkDetailDTO{
		WorkDTO:       ToWorkDTO(work),
		Acts:          actsDTO,
		Collaborators: collabsDTO,
	}
}

// ─── Requests ───────────────────────────────────────────────────────────────

type CreateWorkRequest struct {
	BranchID      string   `json:"branch_id" binding:"required"`
	ClientID      string   `json:"client_id" binding:"required"`
	ActIDs        []string `json:"act_ids" binding:"required,min=1"`
	MainDrafterID *string  `json:"main_drafter_id,omitempty"`
	Folio         *string  `json:"folio,omitempty"`
	Deadline      *string  `json:"deadline,omitempty"`
}

type UpdateWorkRequest struct {
	Folio    *string  `json:"folio,omitempty"`
	Deadline *string  `json:"deadline,omitempty"`
	ActIDs   []string `json:"act_ids,omitempty"`
}

type UpdateWorkStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AddCollaboratorRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type AddCommentRequest struct {
	Message string `json:"message" binding:"required"`
}
