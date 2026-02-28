package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/ports"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/repository"
)

type UpdateEmployeeRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Password  *string `json:"password,omitempty" binding:"omitempty,min=6"`
	Phone     *string `json:"phone,omitempty"`
	Role      *string `json:"role,omitempty"`
	Status    *string `json:"status,omitempty"`
	BranchID  *string `json:"branch_id,omitempty"`
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
}

type UpdateEmployeeUseCase struct {
	repo   repository.UserRepository
	hasher ports.PasswordHasher
}

func NewUpdateEmployeeUseCase(r repository.UserRepository, h ports.PasswordHasher) *UpdateEmployeeUseCase {
	return &UpdateEmployeeUseCase{repo: r, hasher: h}
}

func (uc *UpdateEmployeeUseCase) Execute(ctx context.Context, targetUserID string, requesterRole string, req UpdateEmployeeRequest) error {
	parsedID, err := uuid.Parse(targetUserID)
	if err != nil {
		return errors.New("ID de empleado destino inválido")
	}

	targetUser, err := uc.repo.GetByID(ctx, parsedID)
	if err != nil || targetUser == nil {
		return errors.New("empleado no encontrado")
	}

	targetRole := string(targetUser.Role)
	
	if requesterRole == string(entities.RoleLocalAdmin) {
		if targetRole == string(entities.RoleSuperAdmin) || targetRole == string(entities.RoleLocalAdmin) {
			return errors.New("operación denegada: un administrador no puede modificar a otro administrador o notario")
		}
		if req.Role != nil && *req.Role == string(entities.RoleSuperAdmin) {
			return errors.New("operación denegada: no tienes permisos para asignar el rol de notario titular")
		}
	} else if requesterRole == string(entities.RoleSuperAdmin) {
		if targetRole == string(entities.RoleSuperAdmin) {
			return errors.New("operación denegada: un notario no puede modificar a otro notario titular")
		}
	} else {
		return errors.New("no tienes permisos para usar esta herramienta")
	}

	if req.FullName != nil {
		targetUser.FullName = *req.FullName
	}
	if req.Role != nil {
		targetUser.Role = entities.UserRole(*req.Role)
	}
	if req.Status != nil {
		targetUser.Status = entities.UserStatus(*req.Status)
	}
	if req.Phone != nil {
		targetUser.Phone = req.Phone
	}
	if req.StartTime != nil {
		targetUser.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		targetUser.EndTime = req.EndTime
	}
	if req.BranchID != nil {
		if *req.BranchID == "" {
			targetUser.BranchID = nil
		} else {
			bID, err := uuid.Parse(*req.BranchID)
			if err == nil {
				targetUser.BranchID = &bID
			}
		}
	}

	if req.Email != nil && *req.Email != targetUser.Email {
		existing, _ := uc.repo.GetByEmail(ctx, *req.Email)
		if existing != nil {
			return errors.New("el correo electrónico ya está registrado en otra cuenta")
		}
		targetUser.Email = *req.Email
	}

	if req.Password != nil && *req.Password != "" {
		hashed, err := uc.hasher.HashPassword(*req.Password)
		if err != nil {
			return errors.New("error al procesar la nueva contraseña")
		}
		targetUser.PasswordHash = hashed
	}

	return uc.repo.Update(ctx, targetUser)
}