package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/repository"
	"github.com/gin-gonic/gin"
)

type SearchWorksQuery struct {
	dtos.PaginationRequest
	dtos.DateRangeRequest
	Search   *string `form:"search"`
	Status   *string `form:"status"`
	BranchID *string `form:"branch_id"`
}

type SearchWorksController struct {
	useCase *app.SearchWorksUseCase
}

func NewSearchWorksController(uc *app.SearchWorksUseCase) *SearchWorksController {
	return &SearchWorksController{useCase: uc}
}

func (ctrl *SearchWorksController) Handle(c *gin.Context) {
	// Extraer datos del JWT inyectados por AuthMiddleware
	userRole, _ := c.Get("userRole")
	userID, _ := c.Get("userID")
	branchID, _ := c.Get("branchID")

	var query SearchWorksQuery
	query.Limit = 10
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	filters := repository.WorkFilters{
		Limit:     query.Limit,
		Offset:    query.Offset,
		Search:    query.Search,
		Status:    query.Status,
		BranchID:  query.BranchID,
		StartDate: query.StartDate,
		EndDate:   query.EndDate,
	}

	works, err := ctrl.useCase.Execute(
		c.Request.Context(),
		userRole.(string),
		userID.(string),
		branchID.(string),
		filters,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar trabajos"})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(works),
		Data:  works,
	})
}
