package app

import (
    "time"
    "github.com/google/uuid"
    "github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

type UserPublicDTO struct {
    ID        uuid.UUID `json:"id"`
    FullName  string    `json:"full_name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    Status    string    `json:"status"`
    HireDate  time.Time `json:"hire_date"`
    StartTime *string   `json:"start_time,omitempty"`
    EndTime   *string   `json:"end_time,omitempty"`
}

func ToUserPublicDTO(user *entities.User) UserPublicDTO {
    return UserPublicDTO{
        ID:        user.ID,
        FullName:  user.FullName,
        Email:     user.Email,
        Role:      string(user.Role),
        Status:    string(user.Status),
        HireDate:  user.HireDate,
        StartTime: user.StartTime,
        EndTime:   user.EndTime,
    }
}