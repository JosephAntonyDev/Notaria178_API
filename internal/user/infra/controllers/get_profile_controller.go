package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
)

type GetProfileController struct {
	useCase *app.GetProfileUseCase
}

func NewGetProfileController(uc *app.GetProfileUseCase) *GetProfileController {
	return &GetProfileController{
		useCase: uc,
	}
}

func (ctrl *GetProfileController) Handle(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró información del usuario en el token"})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al procesar el ID de usuario"})
		return
	}

	profile, err := ctrl.useCase.Execute(c.Request.Context(), userIDStr)
	if err != nil {
		if err.Error() == "el usuario no existe en la base de datos" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}