package adapters

import (
	"context"

	auditApp "github.com/JosephAntonyDev/Notaria178_API/internal/audit/app"
	"github.com/google/uuid"
)

// AuditLoggerAdapter adapta *auditApp.LogActionUseCase
// para que cumpla la interfaz work/domain/events.AuditLogger.
type AuditLoggerAdapter struct {
	uc *auditApp.LogActionUseCase
}

func NewAuditLoggerAdapter(uc *auditApp.LogActionUseCase) *AuditLoggerAdapter {
	return &AuditLoggerAdapter{uc: uc}
}

func (a *AuditLoggerAdapter) LogAction(ctx context.Context, action string, entity string, entityID uuid.UUID, userID *uuid.UUID, details interface{}) error {
	var userIDStr *string
	if userID != nil {
		s := userID.String()
		userIDStr = &s
	}
	return a.uc.Execute(ctx, userIDStr, action, entity, entityID.String(), details)
}
