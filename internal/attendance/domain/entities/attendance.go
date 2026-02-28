package entities

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	Date         time.Time  `json:"date"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
}