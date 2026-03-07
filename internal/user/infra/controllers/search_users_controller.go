package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

type SearchUsersQuery struct {
	dtos.PaginationRequest
	Search    *string `form:"search"`
	Status    *string `form:"status"`
	Role      *string `form:"role"`
	BranchID  *string `form:"branch_id"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
	Sort      *string `form:"sort"`
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

	fmt.Printf("[DEBUG] /users/search params | Limit: %d, Offset: %d", query.Limit, query.Offset)
	if query.Search != nil { fmt.Printf(", Search: %s", *query.Search) }
	if query.Status != nil { fmt.Printf(", Status: %s", *query.Status) }
	if query.Role != nil { fmt.Printf(", Role: %s", *query.Role) }
	if query.BranchID != nil { fmt.Printf(", BranchID: %s", *query.BranchID) }
	if query.StartDate != nil { fmt.Printf(", StartDate: %s", *query.StartDate) }
	if query.EndDate != nil { fmt.Printf(", EndDate: %s", *query.EndDate) }
	if query.Sort != nil { fmt.Printf(", Sort: %s", *query.Sort) }
	fmt.Println()

	var filters entities.UserFilters
	filters.Limit = query.Limit
	filters.Offset = query.Offset
	filters.Search = query.Search
	filters.Status = query.Status
	filters.Role = query.Role
	filters.BranchID = query.BranchID
	filters.StartDate = query.StartDate
	filters.EndDate = query.EndDate
	filters.Sort = query.Sort

	users, total, err := ctrl.useCase.Execute(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar usuarios"})
		return
	}

	c.JSON(http.StatusOK, dtos.PaginatedResponse{
		Total: total, 
		Data:  users,
	})
}