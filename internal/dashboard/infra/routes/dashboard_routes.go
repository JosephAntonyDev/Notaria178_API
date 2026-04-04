package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
)

func SetupDashboardRoutes(
	r *gin.Engine,
	ctrl *controllers.DashboardController,
	jwtSecret string,
) {
	api := r.Group("/dashboard")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	// Todos los roles autenticados pueden acceder; el aislamiento de datos
	// se maneja en cada UseCase según el rol del usuario.
	{
		api.GET("/kpis", ctrl.HandleKPIs)
		api.GET("/trend", ctrl.HandleTrend)
		api.GET("/distribution", ctrl.HandleDistribution)
		api.GET("/activity", ctrl.HandleActivity)
		api.GET("/top-drafters", ctrl.HandleTopDrafters)
		api.GET("/top-acts", ctrl.HandleTopActs)
	}
}
