package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/gin-gonic/gin"
)

type GetMyNotificationsQuery struct {
	Page   int   `form:"page"`
	Size   int   `form:"size"`
	IsRead *bool `form:"is_read"`
}

type GetMyNotificationsController struct {
	useCase *app.GetMyNotificationsUseCase
}

func NewGetMyNotificationsController(uc *app.GetMyNotificationsUseCase) *GetMyNotificationsController {
	return &GetMyNotificationsController{useCase: uc}
}

func (ctrl *GetMyNotificationsController) Handle(c *gin.Context) {
	userID, _ := c.Get("userID")

	var query GetMyNotificationsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de consulta inválidos"})
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Size < 1 || query.Size > 50 {
		query.Size = 20
	}

	filters := repository.NotificationFilters{
		Limit:  query.Size,
		Offset: (query.Page - 1) * query.Size,
		IsRead: query.IsRead,
	}

	result, err := ctrl.useCase.Execute(c.Request.Context(), userID.(string), filters)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":         result.Notifications,
		"unread_count": result.UnreadCount,
		"page":         query.Page,
		"size":         query.Size,
	})
}
