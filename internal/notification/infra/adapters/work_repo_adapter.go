package adapters

import (
	"context"

	notifApp "github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/google/uuid"
)

// WorkRepoAdapter adapta el repositorio de Work para que el modulo de
// notificaciones pueda consultar los colaboradores de un trabajo sin
// depender directamente del paquete work/infra/repository.
// Cumple la interfaz notifApp.WorkRepository (solo GetCollaborators).
type WorkRepoAdapter struct {
	// Se usa una interfaz interna para evitar importar el repo concreto
	getter CollaboratorGetter
}

// CollaboratorGetter es la capacidad minima que necesitamos del work repo.
type CollaboratorGetter interface {
	GetCollaborators(ctx context.Context, workID uuid.UUID) ([]entities.WorkCollaboratorInfo, error)
}

func NewWorkRepoAdapter(g CollaboratorGetter) *WorkRepoAdapter {
	return &WorkRepoAdapter{getter: g}
}

// GetCollaborators implementa notifApp.WorkRepository
func (a *WorkRepoAdapter) GetCollaborators(ctx context.Context, workID uuid.UUID) ([]notifApp.WorkCollaboratorInfo, error) {
	collabs, err := a.getter.GetCollaborators(ctx, workID)
	if err != nil {
		return nil, err
	}

	result := make([]notifApp.WorkCollaboratorInfo, len(collabs))
	for i, c := range collabs {
		result[i] = notifApp.WorkCollaboratorInfo{
			UserID:   c.UserID,
			FullName: c.FullName,
		}
	}
	return result, nil
}
