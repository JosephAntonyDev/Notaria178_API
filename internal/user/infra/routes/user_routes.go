package routes

import (
	"github.com/gin-gonic/gin"
	
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

func SetupUserRoutes(
	r *gin.Engine, 
	createUserCtrl *controllers.CreateUserController, 
	loginUserCtrl *controllers.LoginUserController,
	jwtSecret string,
) {
	api := r.Group("/api/v1/users")
	{
		// Rutas Públicas
		api.POST("/login", loginUserCtrl.Handle)

		// Rutas Restringidas a SuperAdmin y LocalAdmin
		adminOnly := api.Group("")
		
		adminOnly.Use(middleware.AuthMiddleware(jwtSecret))
		adminOnly.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
		{
			adminOnly.POST("/create", createUserCtrl.Handle)
		}
	}
}