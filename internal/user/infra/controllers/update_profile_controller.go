package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
)

type UpdateProfileController struct {
	useCase *app.UpdateProfileUseCase
}

func NewUpdateProfileController(uc *app.UpdateProfileUseCase) *UpdateProfileController {
	return &UpdateProfileController{useCase: uc}
}

func (ctrl *UpdateProfileController) Handle(c *gin.Context) {
	var req app.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userIDVal, _ := c.Get("userID")
	
	err := ctrl.useCase.Execute(c.Request.Context(), userIDVal.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado exitosamente"})
}