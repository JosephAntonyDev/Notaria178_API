package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/branch/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func SetupBranchRoutes(
	r *gin.Engine,
	createBranchCtrl *controllers.CreateBranchController,
	updateBranchCtrl *controllers.UpdateBranchController,
	searchBranchesCtrl *controllers.SearchBranchesController,
	jwtSecret string,
) {
	api := r.Group("/branches")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Accesible para cualquier empleado logueado
		api.GET("/search", searchBranchesCtrl.Handle)

		// Restringido a administradores (solo SUPER_ADMIN gestiona sucursales)
		adminOnly := api.Group("")
		adminOnly.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
		{
			adminOnly.POST("/create", createBranchCtrl.Handle)
			adminOnly.PATCH("/update/:id", updateBranchCtrl.Handle)
		}
	}
}
