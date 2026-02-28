package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
)

type CreateUserController struct {
	useCase *app.CreateUserUseCase
}

func NewCreateUserController(uc *app.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{
		useCase: uc,
	}
}

func (ctrl *CreateUserController) Handle(c *gin.Context) {
	requesterRoleVal, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar tu nivel de acceso"})
		return
	}
	requesterRole := requesterRoleVal.(string)

	var req app.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos requeridos"})
		return
	}

	user, err := ctrl.useCase.Execute(c.Request.Context(), requesterRole, req)
	if err != nil {
		if err.Error()[:18] == "operación denegada" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "el correo ya está registrado en el sistema" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado exitosamente",
		"data":    user,
	})
}