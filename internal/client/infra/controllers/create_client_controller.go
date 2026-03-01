package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/app"
	"github.com/gin-gonic/gin"
)

type CreateClientController struct {
	useCase *app.CreateClientUseCase
}

func NewCreateClientController(uc *app.CreateClientUseCase) *CreateClientController {
	return &CreateClientController{useCase: uc}
}

func (ctrl *CreateClientController) Handle(c *gin.Context) {
	var req app.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos requeridos"})
		return
	}

	client, err := ctrl.useCase.Execute(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "el RFC ya está registrado en el sistema" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cliente creado exitosamente",
		"data":    client,
	})
}
