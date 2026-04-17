package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type AddCommentUseCase struct {
	repo             repository.WorkRepository
	audit            events.AuditLogger
	commentNotifier  events.CommentNotifier
}

func NewAddCommentUseCase(r repository.WorkRepository, audit events.AuditLogger, commentNotifier events.CommentNotifier) *AddCommentUseCase {
	return &AddCommentUseCase{repo: r, audit: audit, commentNotifier: commentNotifier}
}

func (uc *AddCommentUseCase) Execute(ctx context.Context, reqCtx RequestContext, workID string, req AddCommentRequest) (*WorkCommentDTO, error) {
	parsedWorkID, err := uuid.Parse(workID)
	if err != nil {
		return nil, errors.New("ID de trabajo invalido")
	}

	work, err := uc.repo.GetByID(ctx, parsedWorkID)
	if err != nil {
		return nil, err
	}
	if work == nil {
		return nil, errors.New("trabajo no encontrado")
	}

	userUUID, err := uuid.Parse(reqCtx.UserID)
	if err != nil {
		return nil, errors.New("error interno al identificar al usuario")
	}

	isCollab, _ := uc.repo.IsCollaborator(ctx, work.ID, userUUID)
	if !CanAccessWork(work, reqCtx, isCollab) {
		return nil, errors.New("no tienes acceso a este trabajo para comentar")
	}

	comment := &entities.WorkComment{
		ID:        uuid.New(),
		WorkID:    parsedWorkID,
		UserID:    userUUID,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.AddComment(ctx, comment); err != nil {
		return nil, err
	}

	// Recuperar el nombre del autor para usarlo en el DTO y notificaciones
	authorName, errName := uc.repo.GetUserFullNameByID(ctx, userUUID)
	if errName != nil {
		fmt.Printf("[WARN] No se pudo obtener el nombre del usuario %s: %v\n", userUUID, errName)
		authorName = "Usuario"
	}
	comment.UserName = authorName

	if uc.audit != nil {
		details := map[string]interface{}{
			"message": req.Message,
		}
		if work.Folio != nil {
			details["folio"] = *work.Folio
		}

		_ = uc.audit.LogAction(ctx, "COMMENT", "WORK", parsedWorkID, &userUUID, details)
	}

	// Disparar notificaciones push + in-app a los colaboradores (en goroutine para no bloquear)
	if uc.commentNotifier != nil {
		folio := "Sin folio"
		if work.Folio != nil {
			folio = *work.Folio
		}

		go func() {
			notifCtx := context.Background()
			err := uc.commentNotifier.NotifyNewComment(notifCtx, events.CommentNotification{
				WorkID:         parsedWorkID,
				WorkFolio:      folio,
				CommentID:      comment.ID,
				CommentAuthor:  userUUID,
				CommentMessage: req.Message,
				AuthorName:     authorName,
			})
			if err != nil {
				fmt.Printf("[ERROR] Fallo al notificar comentario: %v\n", err)
			}
		}()
	}

	dto := ToWorkCommentDTO(*comment)
	return &dto, nil
}
