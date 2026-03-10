package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type DeleteWorkRequirementUseCase struct {
	repo    repository.WorkRepository
	fileDel repository.FileDeleter
	audit   events.AuditLogger
}

func NewDeleteWorkRequirementUseCase(r repository.WorkRepository, fd repository.FileDeleter, audit events.AuditLogger) *DeleteWorkRequirementUseCase {
	return &DeleteWorkRequirementUseCase{repo: r, fileDel: fd, audit: audit}
}

func (uc *DeleteWorkRequirementUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, reqID string) error {
	parsedWorkID, err := uuid.Parse(workID)
	if err != nil {
		return errors.New("ID de trabajo inválido")
	}

	parsedReqID, err := uuid.Parse(reqID)
	if err != nil {
		return errors.New("ID de requisito inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedWorkID)
	if err != nil {
		return err
	}
	if work == nil {
		return errors.New("trabajo no encontrado")
	}

	userUUID, _ := uuid.Parse(reqCtx.UserID)
	isCollab, _ := uc.repo.IsCollaborator(ctx, work.ID, userUUID)

	if !canAccessWork(work, reqCtx, isCollab) {
		return errors.New("no tienes acceso a este trabajo")
	}

	if work.Status == entities.StatusApproved {
		return errors.New("no se puede modificar un trabajo aprobado")
	}

	// Clean up documents linked to this work requirement by requirement_id
	docs, _ := uc.repo.GetDocumentsForCleanupByReqIDs(ctx, parsedWorkID, []uuid.UUID{parsedReqID})
	if len(docs) > 0 {
		docIDs := make([]uuid.UUID, len(docs))
		for i, d := range docs {
			if d.FilePath != "" {
				_ = uc.fileDel.DeleteFile(d.FilePath)
			}
			docIDs[i] = d.ID
		}
		_ = uc.repo.DeleteDocumentRecords(ctx, docIDs)
	}

	// Get work requirement name before deleting
	wrInfo, _ := uc.repo.GetWorkRequirementByID(ctx, parsedReqID)

	if err := uc.repo.DeleteWorkRequirement(ctx, parsedReqID); err != nil {
		return errors.New("error al eliminar el requisito")
	}

	if uc.audit != nil {
		var reqUUID *uuid.UUID
		if parsed, err := uuid.Parse(reqCtx.UserID); err == nil {
			reqUUID = &parsed
		}

		details := map[string]interface{}{}
		if work.Folio != nil {
			details["folio"] = *work.Folio
		}
		if wrInfo != nil {
			details["requirement_name"] = wrInfo.Name
		}

		_ = uc.audit.LogAction(ctx, "DELETE_REQUIREMENT", "WORK", work.ID, reqUUID, details)
	}

	return nil
}
