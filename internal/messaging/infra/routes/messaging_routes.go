package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
)

func SetupMessagingRoutes(r *gin.Engine, wsCtrl *controllers.WSController, jwtSecret string) {
	api := r.Group("/ws")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/comments", wsCtrl.Handle)
	}
}
