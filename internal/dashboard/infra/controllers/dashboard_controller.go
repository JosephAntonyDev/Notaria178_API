package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/app"
	"github.com/gin-gonic/gin"
)

// ─── Query structs para bind de parámetros ──────────────────────────────────

type DashboardQuery struct {
	BranchID  *string `form:"branch_id"`
	Timeframe string  `form:"timeframe"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

type TrendQuery struct {
	DashboardQuery
	GroupBy string `form:"group_by"`
}

type ActivityQuery struct {
	DashboardQuery
	UserID   *string `form:"user_id"`
	EntityID *string `form:"entity_id"`
	Limit    int     `form:"limit"`
	Offset   int     `form:"offset"`
}

type RankingQuery struct {
	DashboardQuery
	Limit int `form:"limit"`
}

// ─── Controller ─────────────────────────────────────────────────────────────

type DashboardController struct {
	getKPIs         *app.GetKPIsUseCase
	getTrend        *app.GetTrendUseCase
	getDistribution *app.GetDistributionUseCase
	getActivity     *app.GetRecentActivityUseCase
	getTopDrafters  *app.GetTopDraftersUseCase
	getTopActs      *app.GetTopActsUseCase
}

func NewDashboardController(
	getKPIs *app.GetKPIsUseCase,
	getTrend *app.GetTrendUseCase,
	getDistribution *app.GetDistributionUseCase,
	getActivity *app.GetRecentActivityUseCase,
	getTopDrafters *app.GetTopDraftersUseCase,
	getTopActs *app.GetTopActsUseCase,
) *DashboardController {
	return &DashboardController{
		getKPIs:         getKPIs,
		getTrend:        getTrend,
		getDistribution: getDistribution,
		getActivity:     getActivity,
		getTopDrafters:  getTopDrafters,
		getTopActs:      getTopActs,
	}
}

// ─── GET /dashboard/kpis ────────────────────────────────────────────────────

func (ctrl *DashboardController) HandleKPIs(c *gin.Context) {
	var query DashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	result, err := ctrl.getKPIs.Execute(
		c.Request.Context(),
		query.BranchID,
		query.Timeframe,
		query.StartDate,
		query.EndDate,
	)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ─── GET /dashboard/trend ───────────────────────────────────────────────────

func (ctrl *DashboardController) HandleTrend(c *gin.Context) {
	var query TrendQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	result, err := ctrl.getTrend.Execute(
		c.Request.Context(),
		query.DashboardQuery.BranchID,
		query.DashboardQuery.Timeframe,
		query.DashboardQuery.StartDate,
		query.DashboardQuery.EndDate,
		query.GroupBy,
	)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ─── GET /dashboard/distribution ────────────────────────────────────────────

func (ctrl *DashboardController) HandleDistribution(c *gin.Context) {
	var query DashboardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	result, err := ctrl.getDistribution.Execute(
		c.Request.Context(),
		query.BranchID,
		query.Timeframe,
		query.StartDate,
		query.EndDate,
	)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ─── GET /dashboard/activity ────────────────────────────────────────────────

func (ctrl *DashboardController) HandleActivity(c *gin.Context) {
	var query ActivityQuery
	query.Limit = 20
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	if query.Limit <= 0 || query.Limit > 100 {
		query.Limit = 20
	}

	result, err := ctrl.getActivity.Execute(
		c.Request.Context(),
		query.DashboardQuery.BranchID,
		query.UserID,
		query.EntityID,
		query.DashboardQuery.Timeframe,
		query.DashboardQuery.StartDate,
		query.DashboardQuery.EndDate,
		query.Limit,
		query.Offset,
	)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ─── GET /dashboard/top-drafters ────────────────────────────────────────────

func (ctrl *DashboardController) HandleTopDrafters(c *gin.Context) {
	var query RankingQuery
	query.Limit = 10

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			query.Limit = parsed
		}
	}

	result, err := ctrl.getTopDrafters.Execute(
		c.Request.Context(),
		query.DashboardQuery.BranchID,
		query.DashboardQuery.Timeframe,
		query.DashboardQuery.StartDate,
		query.DashboardQuery.EndDate,
		query.Limit,
	)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ─── GET /dashboard/top-acts ────────────────────────────────────────────────

func (ctrl *DashboardController) HandleTopActs(c *gin.Context) {
	var query RankingQuery
	query.Limit = 10

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			query.Limit = parsed
		}
	}

	result, err := ctrl.getTopActs.Execute(
		c.Request.Context(),
		query.DashboardQuery.BranchID,
		query.DashboardQuery.Timeframe,
		query.DashboardQuery.StartDate,
		query.DashboardQuery.EndDate,
		query.Limit,
	)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
