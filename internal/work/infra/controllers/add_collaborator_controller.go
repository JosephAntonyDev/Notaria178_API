package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type AddCollaboratorController struct {
	useCase *app.AddCollaboratorUseCase
}

func NewAddCollaboratorController(uc *app.AddCollaboratorUseCase) *AddCollaboratorController {
	return &AddCollaboratorController{useCase: uc}
}

func (ctrl *AddCollaboratorController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	var req app.AddCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: el user_id es requerido"})
		return
	}

	if err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, req); err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Colaborador agregado exitosamente"})
}
