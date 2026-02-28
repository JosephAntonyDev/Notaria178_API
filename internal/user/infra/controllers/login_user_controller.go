package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
)

type LoginUserController struct {
	useCase *app.LoginUserUseCase
}

func NewLoginUserController(uc *app.LoginUserUseCase) *LoginUserController {
	return &LoginUserController{
		useCase: uc,
	}
}

func (ctrl *LoginUserController) Handle(c *gin.Context) {
	var req app.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan datos requeridos o el formato es incorrecto"})
		return
	}

	response, err := ctrl.useCase.Execute(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "credenciales inválidas" || err.Error() == "esta cuenta ha sido desactivada, contacta al administrador" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bienvenido al sistema Notaría 178",
		"data":    response,
	})
}