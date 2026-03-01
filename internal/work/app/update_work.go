package app

import (
	"context"
	"errors"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type UpdateWorkUseCase struct {
	repo repository.WorkRepository
}

func NewUpdateWorkUseCase(r repository.WorkRepository) *UpdateWorkUseCase {
	return &UpdateWorkUseCase{repo: r}
}

func (uc *UpdateWorkUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req UpdateWorkRequest) (*WorkDetailDTO, error) {
	parsedID, err := uuid.Parse(workID)
	if err != nil {
		return nil, errors.New("ID de trabajo inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	if work == nil {
		return nil, errors.New("trabajo no encontrado")
	}

	if !canModifyWork(work, reqCtx) {
		return nil, errors.New("no tienes permisos para modificar este trabajo")
	}

	// Actualizar folio
	if req.Folio != nil {
		work.Folio = req.Folio
	}

	// Actualizar fecha límite
	if req.Deadline != nil && *req.Deadline != "" {
		parsed, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			return nil, errors.New("formato de fecha límite inválido, usa YYYY-MM-DD")
		}
		work.Deadline = &parsed
	}

	if err := uc.repo.Update(ctx, work); err != nil {
		return nil, err
	}

	// Reemplazar actos si se proporcionaron
	if len(req.ActIDs) > 0 {
		actIDs := make([]uuid.UUID, 0, len(req.ActIDs))
		for _, aid := range req.ActIDs {
			parsed, err := uuid.Parse(aid)
			if err != nil {
				return nil, errors.New("uno de los IDs de acto es inválido")
			}
			actIDs = append(actIDs, parsed)
		}
		if err := uc.repo.RemoveAllActs(ctx, work.ID); err != nil {
			return nil, err
		}
		if err := uc.repo.AddActs(ctx, work.ID, actIDs); err != nil {
			return nil, err
		}
	}

	// Construir respuesta
	acts, _ := uc.repo.GetActsByWorkID(ctx, work.ID)
	collabs, _ := uc.repo.GetCollaborators(ctx, work.ID)

	return buildWorkDetail(work, acts, collabs), nil
}
