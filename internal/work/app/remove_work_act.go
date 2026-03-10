package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type RemoveWorkActUseCase struct {
	repo    repository.WorkRepository
	fileDel repository.FileDeleter
	audit   events.AuditLogger
}

func NewRemoveWorkActUseCase(r repository.WorkRepository, fd repository.FileDeleter, audit events.AuditLogger) *RemoveWorkActUseCase {
	return &RemoveWorkActUseCase{repo: r, fileDel: fd, audit: audit}
}

func (uc *RemoveWorkActUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, actID string) error {
	parsedWorkID, err := uuid.Parse(workID)
	if err != nil {
		return errors.New("ID de trabajo inválido")
	}

	parsedActID, err := uuid.Parse(actID)
	if err != nil {
		return errors.New("ID de acto inválido")
	}

	work, err := uc.repo.GetByID(ctx, parsedWorkID)
	if err != nil {
		return err
	}
	if work == nil {
		return errors.New("trabajo no encontrado")
	}

	if !canModifyWork(work, reqCtx) {
		return errors.New("no tienes permisos para modificar este trabajo")
	}

	if work.Status == entities.StatusApproved {
		return errors.New("no se puede modificar un trabajo aprobado")
	}

	// Get the requirements of the act being removed
	removedActReqs, _ := uc.repo.GetRequirementsByActIDs(ctx, []uuid.UUID{parsedActID})

	// Get other acts that remain after removal
	allActs, _ := uc.repo.GetActsByWorkID(ctx, work.ID)
	remainingActIDs := make([]uuid.UUID, 0)
	for _, a := range allActs {
		if a.ActID != parsedActID {
			remainingActIDs = append(remainingActIDs, a.ActID)
		}
	}

	// Get requirements from remaining acts to check for shared ones
	var remainingReqs []entities.ActRequirementInfo
	if len(remainingActIDs) > 0 {
		remainingReqs, _ = uc.repo.GetRequirementsByActIDs(ctx, remainingActIDs)
	}

	// Build a set of requirement names that are still used by remaining acts
	sharedNames := make(map[string]bool)
	for _, r := range remainingReqs {
		sharedNames[r.Name] = true
	}

	// Collect orphaned requirement names (exclusive to the removed act)
	orphanedNames := make([]string, 0)
	for _, r := range removedActReqs {
		if !sharedNames[r.Name] {
			orphanedNames = append(orphanedNames, r.Name)
		}
	}

	// Delete orphaned documents from disk and DB
	if len(orphanedNames) > 0 {
		// Get ALL act_requirement IDs for the orphaned names so we find any linked documents
		allReqIDs, _ := uc.repo.GetActRequirementIDsByNames(ctx, orphanedNames)
		if len(allReqIDs) > 0 {
			docs, _ := uc.repo.GetDocumentsForCleanupByReqIDs(ctx, parsedWorkID, allReqIDs)
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
		}
	}

	// Get act name before removing
	actName, _ := uc.repo.GetActNameByID(ctx, parsedActID)

	// Remove the act from the work
	if err := uc.repo.RemoveAct(ctx, work.ID, parsedActID); err != nil {
		return err
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
		details["act_name"] = actName

		_ = uc.audit.LogAction(ctx, "REMOVE_ACT", "WORK", work.ID, reqUUID, details)
	}

	return nil
}
