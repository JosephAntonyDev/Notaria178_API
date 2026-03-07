package app

import (
	"time"

	branchEntities "github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
	"github.com/google/uuid"
)

type BranchDTO struct {
	Name     string  `json:"name"`
	Address  *string `json:"address,omitempty"`
	ImageURL string  `json:"image_url"`
}

type UserPublicDTO struct {
	ID        uuid.UUID  `json:"id"`
	FullName  string     `json:"full_name"`
	Email     string     `json:"email"`
	Phone     *string    `json:"phone,omitempty"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	HireDate  time.Time  `json:"hire_date"`
	StartTime *string    `json:"start_time,omitempty"`
	EndTime   *string    `json:"end_time,omitempty"`
	Branch    *BranchDTO `json:"branch"`
}

func getBranchImageURL(name string) string {
	switch name {
	case "Tuxtla Gutiérrez":
		return "https://images.unsplash.com/photo-1497366216548-37526070297c?auto=format&fit=crop&q=80&w=2069" // Corporate office
	case "San Fernando":
		return "https://images.unsplash.com/photo-1486406146926-c627a92ad1ab?auto=format&fit=crop&q=80&w=2070" // Modern building
	case "CDMX":
		return "https://images.unsplash.com/photo-1431540015161-1bf339c58da3?auto=format&fit=crop&q=80&w=2070" // Skyscraper vibe
	default:
		return "https://images.unsplash.com/photo-1497366216548-37526070297c?auto=format&fit=crop&q=80&w=2069"
	}
}

func ToUserPublicDTO(user *entities.User, branch *branchEntities.Branch) UserPublicDTO {
	var branchDTO *BranchDTO
	if branch != nil {
		branchDTO = &BranchDTO{
			Name:     branch.Name,
			Address:  branch.Address,
			ImageURL: getBranchImageURL(branch.Name), 
		}
	} else {
		branchDTO = &BranchDTO{
			Name:     "Sucursal no asignada",
			Address:  nil,
			ImageURL: "https://images.unsplash.com/photo-1497366216548-37526070297c?auto=format&fit=crop&q=80&w=2069", 
		}
	}

	return UserPublicDTO{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      string(user.Role),
		Status:    string(user.Status),
		HireDate:  user.HireDate,
		StartTime: user.StartTime,
		EndTime:   user.EndTime,
		Branch:    branchDTO,
	}
}
