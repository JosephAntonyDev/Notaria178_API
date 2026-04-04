package app

import (
	"context"
	"errors"

	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/domain/entities"
	workApp "github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	workRepo "github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/google/uuid"
)

type JoinRoomUseCase struct {
	workRepo workRepo.WorkRepository
}

func NewJoinRoomUseCase(wr workRepo.WorkRepository) *JoinRoomUseCase {
	return &JoinRoomUseCase{workRepo: wr}
}

// Execute valida que el usuario puede unirse a la sala
func (uc *JoinRoomUseCase) Execute(ctx context.Context, reqCtx workApp.RequestContext, room *entities.Room) error {
	if room.Type == entities.RoomTypeWorkComments {
		work, err := uc.workRepo.GetByID(ctx, room.EntityID)
		if err != nil {
			return err
		}
		if work == nil {
			return errors.New("trabajo no encontrado")
		}

		userUUID, err := uuid.Parse(reqCtx.UserID)
		if err != nil {
			return errors.New("ID de usuario inválido")
		}

		isCollab, _ := uc.workRepo.IsCollaborator(ctx, work.ID, userUUID)

		// Reutiliza la lógica de acceso existente
		if !workApp.CanAccessWork(work, reqCtx, isCollab) {
			return errors.New("no tienes acceso a este trabajo")
		}
	}
	// Para RoomTypePrivateChat se agregaría validación de participantes

	return nil
}
