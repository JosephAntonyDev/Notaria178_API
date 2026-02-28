package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
)

type UpdateEmployeeController struct {
	useCase *app.UpdateEmployeeUseCase
}

func NewUpdateEmployeeController(uc *app.UpdateEmployeeUseCase) *UpdateEmployeeController {
	return &UpdateEmployeeController{useCase: uc}
}

func (ctrl *UpdateEmployeeController) Handle(c *gin.Context) {
	targetID := c.Param("id") 

	requesterRoleVal, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo identificar tu nivel de acceso"})
		return
	}
	requesterRole := requesterRoleVal.(string)

	var req app.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos o formato incorrecto"})
		return
	}

	err := ctrl.useCase.Execute(c.Request.Context(), targetID, requesterRole, req)
	if err != nil {
		if err.Error()[:18] == "operación denegada" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Datos del empleado actualizados exitosamente"})
}