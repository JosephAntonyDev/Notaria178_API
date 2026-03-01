package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/gin-gonic/gin"
)

type SearchBranchesQuery struct {
	dtos.PaginationRequest
	Search *string `form:"search"`
}

type SearchBranchesController struct {
	useCase *app.SearchBranchesUseCase
}

func NewSearchBranchesController(uc *app.SearchBranchesUseCase) *SearchBranchesController {
	return &SearchBranchesController{useCase: uc}
}

func (ctrl *SearchBranchesController) Handle(c *gin.Context) {
	var query SearchBranchesQuery

	query.Limit = 10
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	filters := repository.BranchFilters{
		Limit:  query.Limit,
		Offset: query.Offset,
		Search: query.Search,
	}

	branches, err := ctrl.useCase.Execute(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar sucursales"})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(branches),
		Data:  branches,
	})
}
