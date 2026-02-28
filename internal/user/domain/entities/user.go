package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string
type UserStatus string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
	RoleLocalAdmin UserRole = "LOCAL_ADMIN"
	RoleDrafter    UserRole = "DRAFTER"
	RoleDataEntry  UserRole = "DATA_ENTRY"

	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
)

type User struct {
	ID           uuid.UUID   `json:"id"`
	BranchID     *uuid.UUID  `json:"branch_id,omitempty"`
	FullName     string      `json:"full_name"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"`
	Phone        *string     `json:"phone,omitempty"`
	Role         UserRole    `json:"role"`
	Status       UserStatus  `json:"status"`
	HireDate     time.Time   `json:"hire_date"`
	StartTime    *string     `json:"start_time,omitempty"`
	EndTime      *string     `json:"end_time,omitempty"`
	PhotoURL     *string     `json:"photo_url,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
