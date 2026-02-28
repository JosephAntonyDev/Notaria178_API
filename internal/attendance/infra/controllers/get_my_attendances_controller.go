package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/repository"
)

type AttendanceQuery struct {
	dtos.PaginationRequest
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

type GetMyAttendancesController struct {
	useCase *app.GetMyAttendancesUseCase
}

func NewGetMyAttendancesController(uc *app.GetMyAttendancesUseCase) *GetMyAttendancesController {
	return &GetMyAttendancesController{useCase: uc}
}

func (ctrl *GetMyAttendancesController) Handle(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado"})
		return
	}
	userIDStr := userIDVal.(string)

	var query AttendanceQuery
	query.Limit = 10
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de filtrado inválidos"})
		return
	}

	filters := repository.AttendanceFilters{
		Limit:     query.Limit,
		Offset:    query.Offset,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	}

	attendances, err := ctrl.useCase.Execute(c.Request.Context(), userIDStr, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el historial"})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(attendances),
		Data:  attendances,
	})
}