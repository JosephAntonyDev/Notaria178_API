package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteRequirementController struct {
	useCase *app.DeleteRequirementUseCase
}

func NewDeleteRequirementController(uc *app.DeleteRequirementUseCase) *DeleteRequirementController {
	return &DeleteRequirementController{useCase: uc}
}

func (ctrl *DeleteRequirementController) Handle(c *gin.Context) {
	actIDParam := c.Param("id")
	actID, err := uuid.Parse(actIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de acto inválido"})
		return
	}

	reqIDParam := c.Param("req_id")
	reqID, err := uuid.Parse(reqIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de requisito inválido"})
		return
	}

	softDeleted, err := ctrl.useCase.Execute(c.Request.Context(), actID, reqID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar requisito: " + err.Error()})
		return
	}

	if softDeleted {
		c.JSON(http.StatusOK, gin.H{
			"message":      "Requisito desactivado exitosamente",
			"soft_deleted": true,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":      "Requisito eliminado exitosamente",
			"soft_deleted": false,
		})
	}
}
