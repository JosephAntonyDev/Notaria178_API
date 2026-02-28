package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/app"
)

type CheckInOutController struct {
	useCase *app.CheckInOutUseCase
}

func NewCheckInOutController(uc *app.CheckInOutUseCase) *CheckInOutController {
	return &CheckInOutController{useCase: uc}
}

func (ctrl *CheckInOutController) Handle(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró el ID del usuario en el token"})
		return
	}
	userIDStr := userIDVal.(string)

	record, msg, err := ctrl.useCase.Execute(c.Request.Context(), userIDStr)
	if err != nil {
		if err.Error() == "ya has completado tu turno del día de hoy" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al procesar la asistencia: " + err.Error()})
		return
	}

	dto := app.ToAttendanceDTO(record)

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
		"data":    dto,
	})
}