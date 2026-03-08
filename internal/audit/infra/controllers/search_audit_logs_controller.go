package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/gin-gonic/gin"
)

// ─── Query struct para bind de parámetros ───────────────────────────────────

type SearchAuditLogsQuery struct {
	dtos.PaginationRequest
	dtos.DateRangeRequest
	UserID   *string `form:"user_id"`
	Action   *string `form:"action"`
	Entity   *string `form:"entity"`
	EntityID *string `form:"entity_id"`
}

// ─── Controller ─────────────────────────────────────────────────────────────

type SearchAuditLogsController struct {
	useCase *app.SearchAuditLogsUseCase
}

func NewSearchAuditLogsController(uc *app.SearchAuditLogsUseCase) *SearchAuditLogsController {
	return &SearchAuditLogsController{useCase: uc}
}

func (ctrl *SearchAuditLogsController) Handle(c *gin.Context) {
	var query SearchAuditLogsQuery
	query.Limit = 20
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	filters := repository.AuditFilters{
		Limit:     query.Limit,
		Offset:    query.Offset,
		UserID:    query.UserID,
		Action:    query.Action,
		Entity:    query.Entity,
		EntityID:  query.EntityID,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	}

	logs, err := ctrl.useCase.Execute(c.Request.Context(), filters)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(logs),
		Data:  logs,
	})
}
