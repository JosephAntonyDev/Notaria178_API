package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

type SearchUsersQuery struct {
	dtos.PaginationRequest
	Search   *string `form:"search"`
	Status   *string `form:"status"`
	Role     *string `form:"role"`
	BranchID *string `form:"branch_id"`
}

type SearchUsersController struct {
	useCase *app.SearchUsersUseCase
}

func NewSearchUsersController(uc *app.SearchUsersUseCase) *SearchUsersController {
	return &SearchUsersController{useCase: uc}
}

func (ctrl *SearchUsersController) Handle(c *gin.Context) {
	var query SearchUsersQuery
	
	query.Limit = 10
	query.Offset = 0

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	var filters entities.UserFilters
	filters.Limit = query.Limit
	filters.Offset = query.Offset
	filters.Search = query.Search
	filters.Status = query.Status
	filters.Role = query.Role
	filters.BranchID = query.BranchID

	users, err := ctrl.useCase.Execute(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar usuarios"})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: len(users), 
		Data:  users,
	})
}