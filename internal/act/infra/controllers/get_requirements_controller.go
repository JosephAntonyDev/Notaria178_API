package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetRequirementsController struct {
	useCase *app.GetRequirementsUseCase
}

func NewGetRequirementsController(uc *app.GetRequirementsUseCase) *GetRequirementsController {
	return &GetRequirementsController{useCase: uc}
}

func (ctrl *GetRequirementsController) Handle(c *gin.Context) {
	actIDParam := c.Param("id")
	actID, err := uuid.Parse(actIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de acto inválido"})
		return
	}

	reqs, err := ctrl.useCase.Execute(c.Request.Context(), actID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener requisitos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reqs,
	})
}
