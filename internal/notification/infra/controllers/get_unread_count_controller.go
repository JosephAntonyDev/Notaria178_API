package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetUnreadCountController struct {
	useCase *app.GetUnreadCountUseCase
}

func NewGetUnreadCountController(uc *app.GetUnreadCountUseCase) *GetUnreadCountController {
	return &GetUnreadCountController{
		useCase: uc,
	}
}

func (ctrl *GetUnreadCountController) Handle(c *gin.Context) {
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

	count, err := ctrl.useCase.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener contador"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unread_count": count,
	})
}
