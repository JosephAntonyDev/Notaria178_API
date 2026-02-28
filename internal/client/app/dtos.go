package app

import (
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/entities"
	"github.com/google/uuid"
)

type ClientDTO struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	RFC       *string   `json:"rfc,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Email     *string   `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func ToClientDTO(client *entities.Client) ClientDTO {
	return ClientDTO{
		ID:        client.ID,
		FullName:  client.FullName,
		RFC:       client.RFC,
		Phone:     client.Phone,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
	}
}

type CreateClientRequest struct {
	FullName string  `json:"full_name" binding:"required"`
	RFC      *string `json:"rfc,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Email    *string `json:"email,omitempty"`
}

type UpdateClientRequest struct {
	FullName *string `json:"full_name,omitempty"`
	RFC      *string `json:"rfc,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Email    *string `json:"email,omitempty"`
}
