package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	branchRepo := repository.NewPostgresBranchRepository(db)

	createBranchUseCase := app.NewCreateBranchUseCase(branchRepo)
	updateBranchUseCase := app.NewUpdateBranchUseCase(branchRepo)
	searchBranchesUseCase := app.NewSearchBranchesUseCase(branchRepo)

	createBranchCtrl := controllers.NewCreateBranchController(createBranchUseCase)
	updateBranchCtrl := controllers.NewUpdateBranchController(updateBranchUseCase)
	searchBranchesCtrl := controllers.NewSearchBranchesController(searchBranchesUseCase)

	routes.SetupBranchRoutes(r, createBranchCtrl, updateBranchCtrl, searchBranchesCtrl, jwtSecret)
}
