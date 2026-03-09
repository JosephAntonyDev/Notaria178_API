package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func SetupAuditRoutes(
	r *gin.Engine,
	searchCtrl *controllers.SearchAuditLogsController,
	metricsCtrl *controllers.GetAuditMetricsController,
	jwtSecret string,
) {
	api := r.Group("/audit")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	api.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
	{
		api.GET("/search", searchCtrl.Handle)
		api.GET("/metrics", metricsCtrl.Handle)
	}
}
