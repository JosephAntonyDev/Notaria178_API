package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/gin-gonic/gin"
)

// ─── Query struct para bind de parámetros ───────────────────────────────────

type GetAuditMetricsQuery struct {
	dtos.DateRangeRequest
}

// ─── Controller ─────────────────────────────────────────────────────────────

type GetAuditMetricsController struct {
	useCase *app.GetAuditMetricsUseCase
}

func NewGetAuditMetricsController(uc *app.GetAuditMetricsUseCase) *GetAuditMetricsController {
	return &GetAuditMetricsController{useCase: uc}
}

func (ctrl *GetAuditMetricsController) Handle(c *gin.Context) {
	var query GetAuditMetricsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	filters := repository.AuditFilters{
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	}

	metrics, err := ctrl.useCase.Execute(c.Request.Context(), filters)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, metrics)
}
