package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
)

func SetupClientRoutes(
	r *gin.Engine,
	createClientCtrl *controllers.CreateClientController,
	updateClientCtrl *controllers.UpdateClientController,
	searchClientsCtrl *controllers.SearchClientsController,
	jwtSecret string,
) {
	api := r.Group("/api/v1/clients")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/search", searchClientsCtrl.Handle)
		api.POST("/create", createClientCtrl.Handle)
		api.PATCH("/update/:id", updateClientCtrl.Handle)
	}
}
