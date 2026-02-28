package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/gin-gonic/gin"
)

type UpdateActController struct {
	useCase *app.UpdateActUseCase
}

func NewUpdateActController(uc *app.UpdateActUseCase) *UpdateActController {
	return &UpdateActController{useCase: uc}
}

func (ctrl *UpdateActController) Handle(c *gin.Context) {
	actID := c.Param("id")

	var req app.UpdateActRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos enviados"})
		return
	}

	act, err := ctrl.useCase.Execute(c.Request.Context(), actID, req)
	if err != nil {
		switch err.Error() {
		case "ID de acto inválido":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "acto no encontrado":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "el nombre del acto ya está registrado":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Acto actualizado exitosamente",
		"data":    act,
	})
}
