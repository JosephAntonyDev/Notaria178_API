package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterDeviceTokenController struct {
	useCase *app.RegisterDeviceTokenUseCase
}

func NewRegisterDeviceTokenController(uc *app.RegisterDeviceTokenUseCase) *RegisterDeviceTokenController {
	return &RegisterDeviceTokenController{
		useCase: uc,
	}
}

func (ctrl *RegisterDeviceTokenController) Handle(c *gin.Context) {
	// Obtener userID del contexto (middleware de autenticación)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}

	userIDStr := userIDVal.(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	var req app.RegisterDeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := ctrl.useCase.Execute(c.Request.Context(), userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar token FCM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token FCM registrado exitosamente",
	})
}
