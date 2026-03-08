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

// Notifier permite enviar notificaciones a usuarios desde el módulo work.
type Notifier interface {
	SendNotification(ctx context.Context, userID uuid.UUID, workID *uuid.UUID, notifType string, message string) error
}
