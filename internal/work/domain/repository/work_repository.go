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
}
