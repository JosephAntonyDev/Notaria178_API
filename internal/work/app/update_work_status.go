package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type UpdateWorkStatusUseCase struct {
	repo     repository.WorkRepository
	audit    events.AuditLogger
	notifier events.Notifier
}

func NewUpdateWorkStatusUseCase(r repository.WorkRepository, audit events.AuditLogger, notifier events.Notifier) *UpdateWorkStatusUseCase {
	return &UpdateWorkStatusUseCase{repo: r, audit: audit, notifier: notifier}
}

func (uc *UpdateWorkStatusUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req UpdateWorkStatusRequest) (*WorkDTO, error) {
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

	// Verificar acceso básico
	userUUID, _ := uuid.Parse(reqCtx.UserID)
	isCollab, _ := uc.repo.IsCollaborator(ctx, work.ID, userUUID)

	if !canAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a este trabajo")
	}

	newStatus := entities.WorkStatus(req.Status)

	// Validar transición de estado según rol
	if err := validateStatusTransition(work, reqCtx, isCollab, newStatus); err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateStatus(ctx, work.ID, newStatus); err != nil {
		return nil, err
	}

	oldStatus := string(work.Status)

	// ─── Efectos secundarios (fire-and-forget) ──────────────────────────
	// Auditoría
	if uc.audit != nil {
		userUUID, _ := uuid.Parse(reqCtx.UserID)
		details := map[string]string{
			"old_status": oldStatus,
			"new_status": string(newStatus),
		}
		_ = uc.audit.LogAction(ctx, "STATUS_CHANGE", "WORK", work.ID, &userUUID, details)
	}

	// Notificación al proyectista principal
	if uc.notifier != nil && work.MainDrafterID != nil {
		msg := fmt.Sprintf("El expediente %s cambió de %s a %s", work.ID.String(), oldStatus, string(newStatus))
		_ = uc.notifier.SendNotification(ctx, *work.MainDrafterID, &work.ID, "STATUS_CHANGE", msg)
	}

	work.Status = newStatus
	dto := ToWorkDTO(work)
	return &dto, nil
}

// validateStatusTransition aplica las reglas de negocio para cambios de estado
//
// Proyectista/Colaborador:
//
//	IN_PROGRESS → READY_FOR_REVIEW  (terminó el trabajo)
//	PENDING     → IN_PROGRESS       (retoma correcciones)
//
// Admin/Notario:
//
//	READY_FOR_REVIEW → APPROVED      (aprueba)
//	READY_FOR_REVIEW → REJECTED      (rechaza)
//	READY_FOR_REVIEW → PENDING       (devuelve para corrección)
//	+ puede hacer las transiciones de proyectista también
func validateStatusTransition(work *entities.Work, reqCtx RequestContext, isCollab bool, newStatus entities.WorkStatus) error {
	current := work.Status

	// Admin / Notario puede hacer todas las transiciones válidas
	if isAdminOrNotario(reqCtx, work) {
		switch {
		case current == entities.StatusReadyForReview && newStatus == entities.StatusApproved:
			return nil
		case current == entities.StatusReadyForReview && newStatus == entities.StatusRejected:
			return nil
		case current == entities.StatusReadyForReview && newStatus == entities.StatusPending:
			return nil
		case current == entities.StatusInProgress && newStatus == entities.StatusReadyForReview:
			return nil
		case current == entities.StatusPending && newStatus == entities.StatusInProgress:
			return nil
		default:
			return errors.New("transición de estado no permitida")
		}
	}

	// Proyectista / Colaborador tiene transiciones limitadas
	if isDrafterOrCollaborator(work, reqCtx, isCollab) {
		switch {
		case current == entities.StatusInProgress && newStatus == entities.StatusReadyForReview:
			return nil
		case current == entities.StatusPending && newStatus == entities.StatusInProgress:
			return nil
		default:
			return errors.New("transición de estado no permitida para tu rol")
		}
	}

	return errors.New("no tienes permisos para cambiar el estado de este trabajo")
}
