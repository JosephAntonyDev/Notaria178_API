package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/app"
	"github.com/gin-gonic/gin"
)

type UpdateBranchController struct {
	useCase *app.UpdateBranchUseCase
}

func NewUpdateBranchController(uc *app.UpdateBranchUseCase) *UpdateBranchController {
	return &UpdateBranchController{useCase: uc}
}

func (ctrl *UpdateBranchController) Handle(c *gin.Context) {
	branchID := c.Param("id")

	var req app.UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: verifica los campos enviados"})
		return
	}

	branch, err := ctrl.useCase.Execute(c.Request.Context(), branchID, req)
	if err != nil {
		switch err.Error() {
		case "ID de sucursal inválido":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "sucursal no encontrada":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "el nombre de la sucursal ya está registrado":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sucursal actualizada exitosamente",
		"data":    branch,
	})
}
