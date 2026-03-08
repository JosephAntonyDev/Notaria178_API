package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func SetupDashboardRoutes(
	r *gin.Engine,
	ctrl *controllers.DashboardController,
	jwtSecret string,
) {
	api := r.Group("/dashboard")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	api.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
	{
		api.GET("/kpis", ctrl.HandleKPIs)
		api.GET("/trend", ctrl.HandleTrend)
		api.GET("/distribution", ctrl.HandleDistribution)
		api.GET("/activity", ctrl.HandleActivity)
		api.GET("/top-drafters", ctrl.HandleTopDrafters)
		api.GET("/top-acts", ctrl.HandleTopActs)
	}
}
