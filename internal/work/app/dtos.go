package app

import (
	"strings"
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
	Acts             []WorkActInfoDTO      `json:"acts"`
	Collaborators    []WorkCollaboratorDTO `json:"collaborators"`
	Requirements     []DeduplicatedReqDTO  `json:"requirements"`
	WorkRequirements []WorkRequirementDTO  `json:"work_requirements"`
	ClientName       string                `json:"client_name"`
	ClientInfo       *ClientInfoDTO        `json:"client_info,omitempty"`
	BranchName       string                `json:"branch_name"`
	MainDrafterName  string                `json:"main_drafter_name"`
}

type WorkActInfoDTO struct {
	ActID       uuid.UUID `json:"act_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Category    string    `json:"category"`
	Status      string    `json:"status"`
}

type WorkCollaboratorDTO struct {
	UserID   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
}

// ClientInfoDTO contiene los datos completos del cliente para la vista de detalle
type ClientInfoDTO struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	RFC      *string   `json:"rfc,omitempty"`
	Phone    *string   `json:"phone,omitempty"`
	Email    *string   `json:"email,omitempty"`
}

// DeduplicatedReqDTO representa un requisito deduplicado (puede venir de múltiples actos)
type DeduplicatedReqDTO struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	SourceActs []string `json:"source_acts"`
	DocumentID *string  `json:"document_id,omitempty"`
}

// WorkRequirementDTO representa un requisito extra/ad-hoc del trabajo
type WorkRequirementDTO struct {
	ID         uuid.UUID  `json:"id"`
	WorkID     uuid.UUID  `json:"work_id"`
	Name       string     `json:"name"`
	DocumentID *uuid.UUID `json:"document_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
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
	return WorkActInfoDTO{
		ActID:       info.ActID,
		Name:        info.Name,
		Description: info.Description,
		Category:    info.Category,
		Status:      info.Status,
	}
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

func buildWorkDetail(work *entities.Work, acts []entities.WorkActInfo, collabs []entities.WorkCollaboratorInfo, clientName, branchName, mainDrafterName string, clientInfo *entities.ClientInfo, dedupReqs []DeduplicatedReqDTO, workReqs []entities.WorkRequirement) *WorkDetailDTO {
	actsDTO := make([]WorkActInfoDTO, 0, len(acts))
	for _, a := range acts {
		actsDTO = append(actsDTO, ToWorkActInfoDTO(a))
	}
	collabsDTO := make([]WorkCollaboratorDTO, 0, len(collabs))
	for _, c := range collabs {
		collabsDTO = append(collabsDTO, ToWorkCollaboratorDTO(c))
	}

	var clientInfoDTO *ClientInfoDTO
	if clientInfo != nil {
		clientInfoDTO = &ClientInfoDTO{
			ID:       clientInfo.ID,
			FullName: clientInfo.FullName,
			RFC:      clientInfo.RFC,
			Phone:    clientInfo.Phone,
			Email:    clientInfo.Email,
		}
	}

	workReqsDTO := make([]WorkRequirementDTO, 0, len(workReqs))
	for _, wr := range workReqs {
		workReqsDTO = append(workReqsDTO, WorkRequirementDTO{
			ID:         wr.ID,
			WorkID:     wr.WorkID,
			Name:       wr.Name,
			DocumentID: wr.DocumentID,
			CreatedAt:  wr.CreatedAt,
		})
	}

	return &WorkDetailDTO{
		WorkDTO:          ToWorkDTO(work),
		Acts:             actsDTO,
		Collaborators:    collabsDTO,
		Requirements:     dedupReqs,
		WorkRequirements: workReqsDTO,
		ClientName:       clientName,
		ClientInfo:       clientInfoDTO,
		BranchName:       branchName,
		MainDrafterName:  mainDrafterName,
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

type AddWorkActRequest struct {
	ActID string `json:"act_id" binding:"required"`
}

type AddWorkRequirementRequest struct {
	Name string `json:"name" binding:"required"`
}

// DeduplicateRequirements agrupa requisitos de actos por nombre (case-insensitive)
// para que el usuario no vea duplicados cuando múltiples actos exigen lo mismo
func DeduplicateRequirements(reqs []entities.ActRequirementInfo, actMap map[uuid.UUID]string) []DeduplicatedReqDTO {
	type dedupEntry struct {
		firstID    string
		name       string
		status     string
		sourceActs []string
	}

	seen := make(map[string]*dedupEntry)
	order := make([]string, 0)

	for _, r := range reqs {
		key := strings.ToLower(strings.TrimSpace(r.Name))
		actName := actMap[r.ActID]
		if entry, ok := seen[key]; ok {
			// Add source act if not already there
			found := false
			for _, s := range entry.sourceActs {
				if s == actName {
					found = true
					break
				}
			}
			if !found {
				entry.sourceActs = append(entry.sourceActs, actName)
			}
			// Keep ACTIVE if any instance is ACTIVE
			if r.Status == "ACTIVE" {
				entry.status = "ACTIVE"
			}
		} else {
			entry := &dedupEntry{
				firstID:    r.ID.String(),
				name:       r.Name,
				status:     r.Status,
				sourceActs: []string{actName},
			}
			seen[key] = entry
			order = append(order, key)
		}
	}

	result := make([]DeduplicatedReqDTO, 0, len(order))
	for _, key := range order {
		entry := seen[key]
		result = append(result, DeduplicatedReqDTO{
			ID:         entry.firstID,
			Name:       entry.name,
			Status:     entry.status,
			SourceActs: entry.sourceActs,
		})
	}
	return result
}
