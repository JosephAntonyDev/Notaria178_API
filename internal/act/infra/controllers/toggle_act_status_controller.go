package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/gin-gonic/gin"
)

type ToggleActStatusController struct {
	useCase *app.ToggleActStatusUseCase
}

func NewToggleActStatusController(uc *app.ToggleActStatusUseCase) *ToggleActStatusController {
	return &ToggleActStatusController{useCase: uc}
}

func (ctrl *ToggleActStatusController) Handle(c *gin.Context) {
	actID := c.Param("id")

	act, err := ctrl.useCase.Execute(c.Request.Context(), actID)
	if err != nil {
		switch err.Error() {
		case "ID de acto inválido":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "acto no encontrado":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estado del acto actualizado exitosamente",
		"data":    act,
	})
}
