package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type UpdateWorkController struct {
	useCase *app.UpdateWorkUseCase
}

func NewUpdateWorkController(uc *app.UpdateWorkUseCase) *UpdateWorkController {
	return &UpdateWorkController{useCase: uc}
}

func (ctrl *UpdateWorkController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	var req app.UpdateWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos enviados"})
		return
	}

	work, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, req)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trabajo actualizado exitosamente",
		"data":    work,
	})
}
