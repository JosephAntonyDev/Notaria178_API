package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/gin-gonic/gin"
)

type MarkAsReadController struct {
	useCase *app.MarkAsReadUseCase
}

func NewMarkAsReadController(uc *app.MarkAsReadUseCase) *MarkAsReadController {
	return &MarkAsReadController{useCase: uc}
}

func (ctrl *MarkAsReadController) Handle(c *gin.Context) {
	notifID := c.Param("id")
	userID, _ := c.Get("userID")

	if err := ctrl.useCase.Execute(c.Request.Context(), notifID, userID.(string)); err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notificación marcada como leída"})
}
