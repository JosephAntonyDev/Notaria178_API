package app

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/entities"
)

type AttendanceDTO struct {
	ID           uuid.UUID  `json:"id"`
	Date         string     `json:"date"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	TotalHours   string     `json:"total_hours"`
}

func ToAttendanceDTO(att *entities.Attendance) AttendanceDTO {
	var totalHours string

	if att.CheckOutTime != nil {
		duration := att.CheckOutTime.Sub(att.CheckInTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60

		totalHours = fmt.Sprintf("%dh %02dm", hours, minutes)
	} else {
		totalHours = "En curso" 
	}

	return AttendanceDTO{
		ID:           att.ID,
		Date:         att.Date.Format("2006-01-02"),
		CheckInTime:  att.CheckInTime,
		CheckOutTime: att.CheckOutTime,
		TotalHours:   totalHours,
	}
}