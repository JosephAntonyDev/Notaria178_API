package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/storage"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string, audit events.AuditLogger, notifier events.Notifier) {
	workRepo := repository.NewPostgresWorkRepository(db)
	fileStorage := storage.NewLocalFileStorage()

	// Casos de uso
	createWorkUC := app.NewCreateWorkUseCase(workRepo)
	getWorkDetailUC := app.NewGetWorkDetailUseCase(workRepo)
	searchWorksUC := app.NewSearchWorksUseCase(workRepo)
	updateWorkUC := app.NewUpdateWorkUseCase(workRepo)
	updateStatusUC := app.NewUpdateWorkStatusUseCase(workRepo, audit, notifier)
	addCollabUC := app.NewAddCollaboratorUseCase(workRepo)
	removeCollabUC := app.NewRemoveCollaboratorUseCase(workRepo)
	addCommentUC := app.NewAddCommentUseCase(workRepo)
	listCommentsUC := app.NewListCommentsUseCase(workRepo)
	addWorkActUC := app.NewAddWorkActUseCase(workRepo)
	removeWorkActUC := app.NewRemoveWorkActUseCase(workRepo, fileStorage)
	addWorkReqUC := app.NewAddWorkRequirementUseCase(workRepo)
	deleteWorkReqUC := app.NewDeleteWorkRequirementUseCase(workRepo, fileStorage)

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
	addWorkActCtrl := controllers.NewAddWorkActController(addWorkActUC)
	removeWorkActCtrl := controllers.NewRemoveWorkActController(removeWorkActUC)
	addWorkReqCtrl := controllers.NewAddWorkRequirementController(addWorkReqUC)
	deleteWorkReqCtrl := controllers.NewDeleteWorkRequirementController(deleteWorkReqUC)

	routes.SetupWorkRoutes(
		r,
		createWorkCtrl, getWorkDetailCtrl, searchWorksCtrl,
		updateWorkCtrl, updateStatusCtrl,
		addCollabCtrl, removeCollabCtrl,
		addCommentCtrl, getCommentsCtrl,
		addWorkActCtrl, removeWorkActCtrl,
		addWorkReqCtrl, deleteWorkReqCtrl,
		jwtSecret,
	)
}
