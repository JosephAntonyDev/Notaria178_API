package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/client/domain/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/gin-gonic/gin"
)

type SearchClientsQuery struct {
	dtos.PaginationRequest
	Search *string `form:"search"`
}

type SearchClientsController struct {
	useCase *app.SearchClientsUseCase
}

func NewSearchClientsController(uc *app.SearchClientsUseCase) *SearchClientsController {
	return &SearchClientsController{useCase: uc}
}

func (ctrl *SearchClientsController) Handle(c *gin.Context) {
	var query SearchClientsQuery

	query.Limit = 10
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	filters := repository.ClientFilters{
		Limit:  query.Limit,
		Offset: query.Offset,
		Search: query.Search,
	}

	clients, err := ctrl.useCase.Execute(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar clientes"})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(clients),
		Data:  clients,
	})
}
