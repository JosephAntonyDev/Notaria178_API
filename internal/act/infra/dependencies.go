package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/act/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	actRepo := repository.NewPostgresActRepository(db)

	createActUseCase := app.NewCreateActUseCase(actRepo)
	updateActUseCase := app.NewUpdateActUseCase(actRepo)
	toggleActStatusUseCase := app.NewToggleActStatusUseCase(actRepo)
	searchActsUseCase := app.NewSearchActsUseCase(actRepo)

	createActCtrl := controllers.NewCreateActController(createActUseCase)
	updateActCtrl := controllers.NewUpdateActController(updateActUseCase)
	toggleStatusCtrl := controllers.NewToggleActStatusController(toggleActStatusUseCase)
	searchActsCtrl := controllers.NewSearchActsController(searchActsUseCase)

	routes.SetupActRoutes(r, createActCtrl, updateActCtrl, toggleStatusCtrl, searchActsCtrl, jwtSecret)
}
