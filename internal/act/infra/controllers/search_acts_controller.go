package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/gin-gonic/gin"
)

type SearchActsQuery struct {
	dtos.PaginationRequest
	Search *string `form:"search"`
	Status *string `form:"status"`
}

type SearchActsController struct {
	useCase *app.SearchActsUseCase
}

func NewSearchActsController(uc *app.SearchActsUseCase) *SearchActsController {
	return &SearchActsController{useCase: uc}
}

func (ctrl *SearchActsController) Handle(c *gin.Context) {
	var query SearchActsQuery

	query.Limit = 10
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	filters := repository.ActFilters{
		Limit:  query.Limit,
		Offset: query.Offset,
		Search: query.Search,
		Status: query.Status,
	}

	acts, err := ctrl.useCase.Execute(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar actos"})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(acts),
		Data:  acts,
	})
}
