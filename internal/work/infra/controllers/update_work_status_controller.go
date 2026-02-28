package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type UpdateWorkStatusController struct {
	useCase *app.UpdateWorkStatusUseCase
}

func NewUpdateWorkStatusController(uc *app.UpdateWorkStatusUseCase) *UpdateWorkStatusController {
	return &UpdateWorkStatusController{useCase: uc}
}

func (ctrl *UpdateWorkStatusController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	var req app.UpdateWorkStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: el status es requerido"})
		return
	}

	work, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, req)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estado del trabajo actualizado exitosamente",
		"data":    work,
	})
}
