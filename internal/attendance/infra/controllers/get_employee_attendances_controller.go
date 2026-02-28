package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/attendance/domain/repository"
)

type GetEmployeeAttendancesController struct {
	useCase *app.GetEmployeeAttendancesUseCase
}

func NewGetEmployeeAttendancesController(uc *app.GetEmployeeAttendancesUseCase) *GetEmployeeAttendancesController {
	return &GetEmployeeAttendancesController{useCase: uc}
}

func (ctrl *GetEmployeeAttendancesController) Handle(c *gin.Context) {
	targetID := c.Param("id")

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

	attendances, err := ctrl.useCase.Execute(c.Request.Context(), targetID, filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(attendances),
		Data:  attendances,
	})
}