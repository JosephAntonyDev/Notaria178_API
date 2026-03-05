package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/controllers"
)

func SetupWorkRoutes(
	r *gin.Engine,
	createWorkCtrl *controllers.CreateWorkController,
	getWorkDetailCtrl *controllers.GetWorkDetailController,
	searchWorksCtrl *controllers.SearchWorksController,
	updateWorkCtrl *controllers.UpdateWorkController,
	updateStatusCtrl *controllers.UpdateWorkStatusController,
	addCollabCtrl *controllers.AddCollaboratorController,
	removeCollabCtrl *controllers.RemoveCollaboratorController,
	addCommentCtrl *controllers.AddCommentController,
	getCommentsCtrl *controllers.GetCommentsController,
	jwtSecret string,
) {
	api := r.Group("/works")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Búsqueda y detalle — accesible para cualquier usuario autenticado
		api.GET("/search", searchWorksCtrl.Handle)
		api.GET("/:id", getWorkDetailCtrl.Handle)

		// Comentarios — cualquier autenticado (access check en use case)
		api.GET("/:id/comments", getCommentsCtrl.Handle)
		api.POST("/:id/comments", addCommentCtrl.Handle)

		// Creación, modificación, gestión de colaboradores — excluye DATA_ENTRY
		managers := api.Group("")
		managers.Use(middleware.RequireRoles(
			entities.RoleSuperAdmin,
			entities.RoleLocalAdmin,
			entities.RoleDrafter,
		))
		{
			managers.POST("/create", createWorkCtrl.Handle)
			managers.PATCH("/update/:id", updateWorkCtrl.Handle)
			managers.PATCH("/status/:id", updateStatusCtrl.Handle)
			managers.POST("/:id/collaborators", addCollabCtrl.Handle)
			managers.DELETE("/:id/collaborators/:userId", removeCollabCtrl.Handle)
		}
	}
}
