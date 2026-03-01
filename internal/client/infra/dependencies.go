package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/client/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/client/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/client/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/client/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	clientRepo := repository.NewPostgresClientRepository(db)

	createClientUseCase := app.NewCreateClientUseCase(clientRepo)
	updateClientUseCase := app.NewUpdateClientUseCase(clientRepo)
	searchClientsUseCase := app.NewSearchClientsUseCase(clientRepo)

	createClientCtrl := controllers.NewCreateClientController(createClientUseCase)
	updateClientCtrl := controllers.NewUpdateClientController(updateClientUseCase)
	searchClientsCtrl := controllers.NewSearchClientsController(searchClientsUseCase)

	routes.SetupClientRoutes(r, createClientCtrl, updateClientCtrl, searchClientsCtrl, jwtSecret)
}
