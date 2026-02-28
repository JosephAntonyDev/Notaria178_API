package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/repository"
)

type GetEmployeeAttendancesUseCase struct {
	repo repository.AttendanceRepository
}

func NewGetEmployeeAttendancesUseCase(r repository.AttendanceRepository) *GetEmployeeAttendancesUseCase {
	return &GetEmployeeAttendancesUseCase{repo: r}
}

func (uc *GetEmployeeAttendancesUseCase) Execute(ctx context.Context, targetUserID string, filters repository.AttendanceFilters) ([]AttendanceDTO, error) {
	userID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, errors.New("ID de empleado inválido")
	}

	attendances, err := uc.repo.ListByUser(ctx, userID, filters)
	if err != nil {
		return nil, errors.New("error al consultar el historial en la base de datos")
	}

	dtos := make([]AttendanceDTO, 0)
	for _, att := range attendances {
		dtos = append(dtos, ToAttendanceDTO(att))
	}

	return dtos, nil
}