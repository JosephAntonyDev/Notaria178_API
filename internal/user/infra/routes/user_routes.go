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
	getProfileCtrl *controllers.GetProfileController,
	searchUsersCtrl *controllers.SearchUsersController,
	updateProfileCtrl *controllers.UpdateProfileController,
	updateEmployeeCtrl *controllers.UpdateEmployeeController,
	jwtSecret string,
) {
	api := r.Group("/api/v1/users")
	{
		// Públicas
		api.POST("/login", loginUserCtrl.Handle)

		// Protegidas
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			protected.GET("/profile", getProfileCtrl.Handle)
			protected.GET("/search", searchUsersCtrl.Handle)
			protected.PATCH("/profile", updateProfileCtrl.Handle)

			// Restringidas
			adminOnly := protected.Group("")
			adminOnly.Use(middleware.RequireRoles(entities.RoleSuperAdmin, entities.RoleLocalAdmin))
			{
				adminOnly.POST("/create", createUserCtrl.Handle)
				adminOnly.PATCH("/update/:id", updateEmployeeCtrl.Handle)
			}
		}
	}
}