package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/act/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func SetupActRoutes(
	r *gin.Engine,
	createActCtrl *controllers.CreateActController,
	updateActCtrl *controllers.UpdateActController,
	toggleStatusCtrl *controllers.ToggleActStatusController,
	searchActsCtrl *controllers.SearchActsController,
	deleteActCtrl *controllers.DeleteActController,
	addReqCtrl *controllers.AddRequirementController,
	delReqCtrl *controllers.DeleteRequirementController,
	getReqsCtrl *controllers.GetRequirementsController,
	jwtSecret string,
) {
	api := r.Group("/acts")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Accesible para cualquier empleado logueado
		api.GET("/search", searchActsCtrl.Handle)
		api.GET("/:id/requirements", getReqsCtrl.Handle)

		// Restringido a administradores
		adminOnly := api.Group("")
		adminOnly.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
		{
			adminOnly.POST("/create", createActCtrl.Handle)
			adminOnly.PATCH("/update/:id", updateActCtrl.Handle)
			adminOnly.PATCH("/status/:id", toggleStatusCtrl.Handle)
			adminOnly.DELETE("/:id", deleteActCtrl.Handle)
			adminOnly.POST("/:id/requirements", addReqCtrl.Handle)
			adminOnly.DELETE("/:id/requirements/:req_id", delReqCtrl.Handle)
		}
	}
}
