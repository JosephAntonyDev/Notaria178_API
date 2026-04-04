package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/controllers"
)

func SetupNotificationRoutes(
	r *gin.Engine,
	getMyNotificationsCtrl *controllers.GetMyNotificationsController,
	markAsReadCtrl *controllers.MarkAsReadController,
	markAllReadCtrl *controllers.MarkAllReadController,
	sseCtrl *controllers.StreamNotificationsController,
	registerDeviceTokenCtrl *controllers.RegisterDeviceTokenController,
	getUnreadCountCtrl *controllers.GetUnreadCountController,
	jwtSecret string,
) {
	api := r.Group("/notifications")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("", getMyNotificationsCtrl.Handle)
		api.GET("/stream", sseCtrl.Handle)
		api.GET("/unread-count", getUnreadCountCtrl.Handle)
		api.PATCH("/:id/read", markAsReadCtrl.Handle)
		api.PATCH("/read-all", markAllReadCtrl.Handle)
		api.PUT("/device-token", registerDeviceTokenCtrl.Handle)
	}
}
