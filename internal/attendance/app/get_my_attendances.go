package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/repository"
)

type GetMyAttendancesUseCase struct {
	repo repository.AttendanceRepository
}

func NewGetMyAttendancesUseCase(r repository.AttendanceRepository) *GetMyAttendancesUseCase {
	return &GetMyAttendancesUseCase{repo: r}
}

func (uc *GetMyAttendancesUseCase) Execute(ctx context.Context, userIDStr string, filters repository.AttendanceFilters) ([]AttendanceDTO, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	attendances, err := uc.repo.ListByUser(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	dtos := make([]AttendanceDTO, 0)
	for _, att := range attendances {
		dtos = append(dtos, ToAttendanceDTO(att))
	}

	return dtos, nil
}