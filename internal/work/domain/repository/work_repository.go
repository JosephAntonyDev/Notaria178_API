package repository

import (
	"context"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

type WorkFilters struct {
	Limit        int
	Offset       int
	Search       *string // búsqueda por folio
	Status       *string
	BranchID     *string // filtro por sucursal (admins, data_entry)
	ScopedUserID *string // solo trabajos donde el usuario es proyectista o colaborador
	StartDate    *string // filtro created_at >= (formato YYYY-MM-DD)
	EndDate      *string // filtro created_at <= (formato YYYY-MM-DD)
	Sort         *string // sorting
}

type WorkRepository interface {
	// CRUD del expediente
	Create(ctx context.Context, work *entities.Work) error
	Update(ctx context.Context, work *entities.Work) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.WorkStatus) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Work, error)
	List(ctx context.Context, filters WorkFilters) ([]*entities.Work, error)

	// Actos del expediente
	AddActs(ctx context.Context, workID uuid.UUID, actIDs []uuid.UUID) error
	RemoveAct(ctx context.Context, workID uuid.UUID, actID uuid.UUID) error
	RemoveAllActs(ctx context.Context, workID uuid.UUID) error
	GetActsByWorkID(ctx context.Context, workID uuid.UUID) ([]entities.WorkActInfo, error)

	// Colaboradores
	AddCollaborator(ctx context.Context, workID uuid.UUID, userID uuid.UUID) error
	RemoveCollaborator(ctx context.Context, workID uuid.UUID, userID uuid.UUID) error
	GetCollaborators(ctx context.Context, workID uuid.UUID) ([]entities.WorkCollaboratorInfo, error)
	IsCollaborator(ctx context.Context, workID uuid.UUID, userID uuid.UUID) (bool, error)

	// Comentarios
	AddComment(ctx context.Context, comment *entities.WorkComment) error
	GetCommentsByWorkID(ctx context.Context, workID uuid.UUID) ([]entities.WorkComment, error)

	// Lookups auxiliares para enriquecer el detalle
	GetClientNameByID(ctx context.Context, clientID uuid.UUID) (string, error)
	GetClientByID(ctx context.Context, clientID uuid.UUID) (*entities.ClientInfo, error)
	GetBranchNameByID(ctx context.Context, branchID uuid.UUID) (string, error)
	GetUserFullNameByID(ctx context.Context, userID uuid.UUID) (string, error)
	UpdateWorkClientID(ctx context.Context, workID uuid.UUID, newClientID uuid.UUID) error
	CountWorksWithClientInStatus(ctx context.Context, clientID uuid.UUID, status string) (int, error)
	GetRequirementsByActIDs(ctx context.Context, actIDs []uuid.UUID) ([]entities.ActRequirementInfo, error)

	// Requisitos extra del trabajo (ad-hoc)
	AddWorkRequirement(ctx context.Context, workID uuid.UUID, name string) (*entities.WorkRequirement, error)
	GetWorkRequirements(ctx context.Context, workID uuid.UUID) ([]entities.WorkRequirement, error)
	DeleteWorkRequirement(ctx context.Context, reqID uuid.UUID) error

	// Lookup de documentos de requisitos subidos para un trabajo (requirement_id → document_id)
	GetRequirementDocumentsByWorkID(ctx context.Context, workID uuid.UUID) (map[uuid.UUID]uuid.UUID, error)

	GetWorkRequirementByID(ctx context.Context, reqID uuid.UUID) (*entities.WorkRequirement, error)
	GetDocumentsForCleanupByReqIDs(ctx context.Context, workID uuid.UUID, reqIDs []uuid.UUID) ([]entities.DocCleanupInfo, error)
	GetActRequirementIDsByNames(ctx context.Context, names []string) ([]uuid.UUID, error)
	DeleteDocumentRecords(ctx context.Context, docIDs []uuid.UUID) error
}

// FileDeleter abstrae el borrado de archivos del disco
type FileDeleter interface {
	DeleteFile(filePath string) error
}
