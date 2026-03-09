package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type AddWorkRequirementController struct {
	useCase *app.AddWorkRequirementUseCase
}

func NewAddWorkRequirementController(uc *app.AddWorkRequirementUseCase) *AddWorkRequirementController {
	return &AddWorkRequirementController{useCase: uc}
}

func (ctrl *AddWorkRequirementController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	var req app.AddWorkRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: se requiere name"})
		return
	}

	result, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, req)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Requisito agregado exitosamente",
		"data":    result,
	})
}
