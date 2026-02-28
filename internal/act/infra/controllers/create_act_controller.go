package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/gin-gonic/gin"
)

type CreateActController struct {
	useCase *app.CreateActUseCase
}

func NewCreateActController(uc *app.CreateActUseCase) *CreateActController {
	return &CreateActController{useCase: uc}
}

func (ctrl *CreateActController) Handle(c *gin.Context) {
	var req app.CreateActRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos requeridos"})
		return
	}

	act, err := ctrl.useCase.Execute(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "el nombre del acto ya está registrado" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Acto creado exitosamente",
		"data":    act,
	})
}
