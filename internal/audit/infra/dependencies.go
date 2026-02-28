package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/audit/infra/routes"
)

// SetupDependencies conecta todas las capas del módulo de auditoría.
// Devuelve el *LogActionUseCase para que main.go pueda inyectarlo
// en los módulos que necesiten registrar acciones de auditoría.
func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) *app.LogActionUseCase {
	auditRepo := repository.NewPostgresAuditRepository(db)

	// Casos de uso
	logActionUC := app.NewLogActionUseCase(auditRepo)
	searchAuditLogsUC := app.NewSearchAuditLogsUseCase(auditRepo)

	// Controladores
	searchCtrl := controllers.NewSearchAuditLogsController(searchAuditLogsUC)

	routes.SetupAuditRoutes(r, searchCtrl, jwtSecret)

	return logActionUC
}
