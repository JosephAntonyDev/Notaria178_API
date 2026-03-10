package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/core/cache"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/storage"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string, audit events.AuditLogger, notifier events.Notifier, cachePort cache.CachePort) {
	workRepo := repository.NewPostgresWorkRepository(db)
	fileStorage := storage.NewLocalFileStorage()

	// Casos de uso
	createWorkUC := app.NewCreateWorkUseCase(workRepo, cachePort, audit)
	getWorkDetailUC := app.NewGetWorkDetailUseCase(workRepo)
	searchWorksUC := app.NewSearchWorksUseCase(workRepo)
	updateWorkUC := app.NewUpdateWorkUseCase(workRepo, audit)
	updateStatusUC := app.NewUpdateWorkStatusUseCase(workRepo, audit, notifier, cachePort)
	addCollabUC := app.NewAddCollaboratorUseCase(workRepo, audit)
	removeCollabUC := app.NewRemoveCollaboratorUseCase(workRepo, audit)
	addCommentUC := app.NewAddCommentUseCase(workRepo, audit)
	listCommentsUC := app.NewListCommentsUseCase(workRepo)
	addWorkActUC := app.NewAddWorkActUseCase(workRepo, audit)
	removeWorkActUC := app.NewRemoveWorkActUseCase(workRepo, fileStorage, audit)
	addWorkReqUC := app.NewAddWorkRequirementUseCase(workRepo, audit)
	deleteWorkReqUC := app.NewDeleteWorkRequirementUseCase(workRepo, fileStorage, audit)

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
