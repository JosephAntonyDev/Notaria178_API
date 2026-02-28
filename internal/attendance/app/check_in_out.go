package app

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/repository"
)

type CheckInOutUseCase struct {
	repo repository.AttendanceRepository
}

func NewCheckInOutUseCase(r repository.AttendanceRepository) *CheckInOutUseCase {
	return &CheckInOutUseCase{repo: r}
}

func (uc *CheckInOutUseCase) Execute(ctx context.Context, userIDStr string) (*entities.Attendance, string, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, "", errors.New("ID de usuario inválido")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	record, err := uc.repo.GetByDateAndUser(ctx, userID, today)
	if err != nil {
		return nil, "", err
	}

	if record == nil {
		newRecord := &entities.Attendance{
			UserID:      userID,
			Date:        today,
			CheckInTime: now,
		}
		err := uc.repo.CreateCheckIn(ctx, newRecord)
		return newRecord, "Entrada registrada exitosamente", err
	}

	if record.CheckOutTime == nil {
		err := uc.repo.UpdateCheckOut(ctx, record.ID, now)
		record.CheckOutTime = &now
		return record, "Salida registrada exitosamente", err
	}

	return record, "", errors.New("ya has completado tu turno del día de hoy")
}