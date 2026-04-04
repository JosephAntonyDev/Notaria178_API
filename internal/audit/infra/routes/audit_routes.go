package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
)

func SetupAuditRoutes(
	r *gin.Engine,
	searchCtrl *controllers.SearchAuditLogsController,
	metricsCtrl *controllers.GetAuditMetricsController,
	jwtSecret string,
) {
	api := r.Group("/audit")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	// Todos los roles autenticados pueden acceder; el aislamiento de datos
	// se maneja en el controller según el rol del usuario.
	{
		api.GET("/search", searchCtrl.Handle)
		api.GET("/metrics", metricsCtrl.Handle)
	}
}
