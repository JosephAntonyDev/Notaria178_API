package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddRequirementController struct {
	useCase *app.AddRequirementUseCase
}

func NewAddRequirementController(uc *app.AddRequirementUseCase) *AddRequirementController {
	return &AddRequirementController{useCase: uc}
}

func (ctrl *AddRequirementController) Handle(c *gin.Context) {
	actIDParam := c.Param("id")
	actID, err := uuid.Parse(actIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de acto inválido"})
		return
	}

	var req app.AddRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de requisito inválido"})
		return
	}
	req.ActID = actID

	result, err := ctrl.useCase.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al agregar requisito: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Requisito agregado exitosamente",
		"data":    result,
	})
}
