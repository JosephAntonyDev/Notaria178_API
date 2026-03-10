package ports

import (
	"context"

	"github.com/google/uuid"
)

// AuditLogger permite registrar acciones de auditoría desde el módulo user.
type AuditLogger interface {
	LogAction(ctx context.Context, action string, entity string, entityID uuid.UUID, userID *uuid.UUID, details interface{}) error
}
