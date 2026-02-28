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
	var req app.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos requeridos"})
		return
	}

	user, err := ctrl.useCase.Execute(c.Request.Context(), req)
	if err != nil {
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