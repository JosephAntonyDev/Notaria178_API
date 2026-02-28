package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type CreateWorkController struct {
	useCase *app.CreateWorkUseCase
}

func NewCreateWorkController(uc *app.CreateWorkUseCase) *CreateWorkController {
	return &CreateWorkController{useCase: uc}
}

func (ctrl *CreateWorkController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)

	var req app.CreateWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos requeridos"})
		return
	}

	work, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, req)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Trabajo creado exitosamente",
		"data":    work,
	})
}
