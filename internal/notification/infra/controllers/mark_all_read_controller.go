package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/gin-gonic/gin"
)

type MarkAllReadController struct {
	useCase *app.MarkAllReadUseCase
}

func NewMarkAllReadController(uc *app.MarkAllReadUseCase) *MarkAllReadController {
	return &MarkAllReadController{useCase: uc}
}

func (ctrl *MarkAllReadController) Handle(c *gin.Context) {
	userID, _ := c.Get("userID")

	if err := ctrl.useCase.Execute(c.Request.Context(), userID.(string)); err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todas las notificaciones marcadas como leídas"})
}
