package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/dashboard/infra/routes"
)

// SetupDependencies conecta todas las capas del módulo de dashboard.
func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string, cachePort cache.CachePort) {
	dashRepo := repository.NewPostgresDashboardRepository(db)

	// Casos de uso
	getKPIsUC := app.NewGetKPIsUseCase(dashRepo, cachePort)
	getTrendUC := app.NewGetTrendUseCase(dashRepo, cachePort)
	getDistributionUC := app.NewGetDistributionUseCase(dashRepo, cachePort)
	getActivityUC := app.NewGetRecentActivityUseCase(dashRepo, cachePort)
	getTopDraftersUC := app.NewGetTopDraftersUseCase(dashRepo, cachePort)
	getTopActsUC := app.NewGetTopActsUseCase(dashRepo, cachePort)

	// Controlador
	ctrl := controllers.NewDashboardController(
		getKPIsUC, getTrendUC, getDistributionUC,
		getActivityUC, getTopDraftersUC, getTopActsUC,
	)

	routes.SetupDashboardRoutes(r, ctrl, jwtSecret)
}
