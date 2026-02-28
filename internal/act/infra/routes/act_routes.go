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
	jwtSecret string,
) {
	api := r.Group("/api/v1/acts")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Accesible para cualquier empleado logueado
		api.GET("/search", searchActsCtrl.Handle)

		// Restringido a administradores
		adminOnly := api.Group("")
		adminOnly.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
		{
			adminOnly.POST("/create", createActCtrl.Handle)
			adminOnly.PATCH("/update/:id", updateActCtrl.Handle)
			adminOnly.PATCH("/status/:id", toggleStatusCtrl.Handle)
		}
	}
}
