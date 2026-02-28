package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/app"
	"github.com/gin-gonic/gin"
)

type UpdateClientController struct {
	useCase *app.UpdateClientUseCase
}

func NewUpdateClientController(uc *app.UpdateClientUseCase) *UpdateClientController {
	return &UpdateClientController{useCase: uc}
}

func (ctrl *UpdateClientController) Handle(c *gin.Context) {
	clientID := c.Param("id")

	var req app.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos enviados"})
		return
	}

	client, err := ctrl.useCase.Execute(c.Request.Context(), clientID, req)
	if err != nil {
		switch err.Error() {
		case "ID de cliente inválido":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "cliente no encontrado":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "el RFC ya está registrado en el sistema":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cliente actualizado exitosamente",
		"data":    client,
	})
}
