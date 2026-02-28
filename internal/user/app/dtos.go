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
	PhotoURL  *string   `json:"photo_url,omitempty"`
	HireDate  time.Time `json:"hire_date"`
}

// Función "Mapeador" (Mapper) para convertir fácilmente una Entidad a un DTO
func ToUserPublicDTO(user *entities.User) UserPublicDTO {
	return UserPublicDTO{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     string(user.Role),
		PhotoURL: user.PhotoURL,
		HireDate: user.HireDate,
	}
}