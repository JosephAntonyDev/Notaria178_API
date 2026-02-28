package app

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/repository"
	"github.com/google/uuid"
)

// LogActionUseCase es un caso de uso INTERNO para registrar acciones de auditoría.
// Será invocado por otros módulos de la API; NO tiene controlador HTTP.
type LogActionUseCase struct {
	repo repository.AuditRepository
}

func NewLogActionUseCase(r repository.AuditRepository) *LogActionUseCase {
	return &LogActionUseCase{repo: r}
}

func (uc *LogActionUseCase) Execute(ctx context.Context, userID *string, action string, entity string, entityID string, details interface{}) error {
	parsedEntityID, err := uuid.Parse(entityID)
	if err != nil {
		return errors.New("ID de entidad inválido")
	}

	var parsedUserID *uuid.UUID
	if userID != nil && *userID != "" {
		uid, err := uuid.Parse(*userID)
		if err != nil {
			return errors.New("ID de usuario inválido")
		}
		parsedUserID = &uid
	}

	var jsonDetails json.RawMessage
	if details != nil {
		raw, err := json.Marshal(details)
		if err != nil {
			return errors.New("error al serializar los detalles de auditoría")
		}
		jsonDetails = raw
	}

	log := &entities.AuditLog{
		ID:          uuid.New(),
		UserID:      parsedUserID,
		Action:      action,
		Entity:      entity,
		EntityID:    parsedEntityID,
		JSONDetails: jsonDetails,
		CreatedAt:   time.Now(),
	}

	return uc.repo.Create(ctx, log)
}
