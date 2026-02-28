package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	workRepo := repository.NewPostgresWorkRepository(db)

	// Casos de uso
	createWorkUC := app.NewCreateWorkUseCase(workRepo)
	getWorkDetailUC := app.NewGetWorkDetailUseCase(workRepo)
	searchWorksUC := app.NewSearchWorksUseCase(workRepo)
	updateWorkUC := app.NewUpdateWorkUseCase(workRepo)
	updateStatusUC := app.NewUpdateWorkStatusUseCase(workRepo)
	addCollabUC := app.NewAddCollaboratorUseCase(workRepo)
	removeCollabUC := app.NewRemoveCollaboratorUseCase(workRepo)
	addCommentUC := app.NewAddCommentUseCase(workRepo)
	listCommentsUC := app.NewListCommentsUseCase(workRepo)

	// Controladores
	createWorkCtrl := controllers.NewCreateWorkController(createWorkUC)
	getWorkDetailCtrl := controllers.NewGetWorkDetailController(getWorkDetailUC)
	searchWorksCtrl := controllers.NewSearchWorksController(searchWorksUC)
	updateWorkCtrl := controllers.NewUpdateWorkController(updateWorkUC)
	updateStatusCtrl := controllers.NewUpdateWorkStatusController(updateStatusUC)
	addCollabCtrl := controllers.NewAddCollaboratorController(addCollabUC)
	removeCollabCtrl := controllers.NewRemoveCollaboratorController(removeCollabUC)
	addCommentCtrl := controllers.NewAddCommentController(addCommentUC)
	getCommentsCtrl := controllers.NewGetCommentsController(listCommentsUC)

	routes.SetupWorkRoutes(
		r,
		createWorkCtrl, getWorkDetailCtrl, searchWorksCtrl,
		updateWorkCtrl, updateStatusCtrl,
		addCollabCtrl, removeCollabCtrl,
		addCommentCtrl, getCommentsCtrl,
		jwtSecret,
	)
}
