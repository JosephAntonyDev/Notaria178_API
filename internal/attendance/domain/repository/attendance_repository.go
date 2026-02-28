package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/entities"
)

type AttendanceFilters struct {
	Limit     int
	Offset    int
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

type AttendanceRepository interface {
	CreateCheckIn(ctx context.Context, attendance *entities.Attendance) error
	UpdateCheckOut(ctx context.Context, id uuid.UUID, checkOutTime time.Time) error
	GetByDateAndUser(ctx context.Context, userID uuid.UUID, date time.Time) (*entities.Attendance, error)
	ListByUser(ctx context.Context, userID uuid.UUID, filters AttendanceFilters) ([]*entities.Attendance, error)
}