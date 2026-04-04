package events

import (
	"context"

	"github.com/google/uuid"
)

// ─── Puertos de Salida (Output Ports) ───────────────────────────────────────
// Estas interfaces permiten que el módulo "work" dispare efectos secundarios
// (auditoría, notificaciones) sin depender directamente de otros módulos.
// Los adaptadores se inyectan desde main.go.

// AuditLogger permite registrar acciones de auditoría desde el módulo work.
type AuditLogger interface {
	LogAction(ctx context.Context, action string, entity string, entityID uuid.UUID, userID *uuid.UUID, details interface{}) error
}

// Notifier permite enviar notificaciones a usuarios desde el modulo work.
type Notifier interface {
	SendNotification(ctx context.Context, userID uuid.UUID, workID *uuid.UUID, notifType string, message string) error
	NotifySuperAdmins(ctx context.Context, workID *uuid.UUID, notifType string, message string) error
}

// CommentNotifier permite disparar notificaciones push + in-app al crear un comentario.
type CommentNotifier interface {
	NotifyNewComment(ctx context.Context, input CommentNotification) error
}

// CommentNotification contiene los datos necesarios para notificar un nuevo comentario.
type CommentNotification struct {
	WorkID         uuid.UUID
	WorkFolio      string
	CommentID      uuid.UUID
	CommentAuthor  uuid.UUID
	CommentMessage string
	AuthorName     string
}
