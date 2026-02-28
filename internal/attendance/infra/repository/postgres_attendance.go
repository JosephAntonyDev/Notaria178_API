package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"strconv"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/repository"

)

type PostgresAttendanceRepository struct {
	db *sql.DB
}

func NewPostgresAttendanceRepository(db *sql.DB) *PostgresAttendanceRepository {
	return &PostgresAttendanceRepository{db: db}
}

func (repo *PostgresAttendanceRepository) CreateCheckIn(ctx context.Context, attendance *entities.Attendance) error {
	query := `
		INSERT INTO attendances (user_id, date, check_in_time) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`
	
	err := repo.db.QueryRowContext(ctx, query,
		attendance.UserID,
		attendance.Date,
		attendance.CheckInTime,
	).Scan(&attendance.ID)

	return err
}

func (repo *PostgresAttendanceRepository) UpdateCheckOut(ctx context.Context, id uuid.UUID, checkOutTime time.Time) error {
	query := `
		UPDATE attendances 
		SET check_out_time = $1 
		WHERE id = $2
	`
	_, err := repo.db.ExecContext(ctx, query, checkOutTime, id)
	return err
}

func (repo *PostgresAttendanceRepository) GetByDateAndUser(ctx context.Context, userID uuid.UUID, date time.Time) (*entities.Attendance, error) {
	query := `
		SELECT id, user_id, date, check_in_time, check_out_time 
		FROM attendances 
		WHERE user_id = $1 AND date = $2
	`
	
	row := repo.db.QueryRowContext(ctx, query, userID, date)

	var attendance entities.Attendance
	err := row.Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.Date,
		&attendance.CheckInTime,
		&attendance.CheckOutTime,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &attendance, nil
}

func (repo *PostgresAttendanceRepository) ListByUser(ctx context.Context, userID uuid.UUID, filters repository.AttendanceFilters) ([]*entities.Attendance, error) {
	baseQuery := `
		SELECT id, user_id, date, check_in_time, check_out_time 
		FROM attendances 
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argId := 2

	if filters.StartDate != nil && *filters.StartDate != "" {
		baseQuery += ` AND date >= $` + strconv.Itoa(argId)
		args = append(args, *filters.StartDate)
		argId++
	}

	if filters.EndDate != nil && *filters.EndDate != "" {
		baseQuery += ` AND date <= $` + strconv.Itoa(argId)
		args = append(args, *filters.EndDate)
		argId++
	}

	baseQuery += ` ORDER BY date DESC LIMIT $` + strconv.Itoa(argId) + ` OFFSET $` + strconv.Itoa(argId+1)
	args = append(args, filters.Limit, filters.Offset)

	rows, err := repo.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*entities.Attendance
	for rows.Next() {
		var att entities.Attendance
		err := rows.Scan(&att.ID, &att.UserID, &att.Date, &att.CheckInTime, &att.CheckOutTime)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, &att)
	}
	return attendances, nil
}