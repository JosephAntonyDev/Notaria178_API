package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/app"
	"github.com/gin-gonic/gin"
)

type CreateBranchController struct {
	useCase *app.CreateBranchUseCase
}

func NewCreateBranchController(uc *app.CreateBranchUseCase) *CreateBranchController {
	return &CreateBranchController{useCase: uc}
}

func (ctrl *CreateBranchController) Handle(c *gin.Context) {
	var req app.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos requeridos"})
		return
	}

	branch, err := ctrl.useCase.Execute(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "el nombre de la sucursal ya está registrado" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sucursal creada exitosamente",
		"data":    branch,
	})
}
