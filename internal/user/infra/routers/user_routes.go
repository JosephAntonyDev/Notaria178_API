package routers

import (
	"github.com/gin-gonic/gin"
	
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/controllers"
	// "github.com/JosephAntonyDev/Notaria178_API/internal/middleware" // Descomentar cuando agreguemos los middlewares
)

func SetupUserRoutes(
	r *gin.Engine, 
	createUserCtrl *controllers.CreateUserController, 
	loginUserCtrl *controllers.LoginUserController,
	// jwtSecret string, // Lo usaremos para el middleware de Auth más adelante
) {
	api := r.Group("/api/v1/users")
	{
		// Rutas Públicas
		api.POST("/login", loginUserCtrl.Handle)

		// Rutas Restringidas (Temporalmente públicas para que puedas probarlas en Postman hoy mismo)
		// Más adelante las envolveremos en: api.Use(middleware.RequireRoles(entities.RoleSuperAdmin))
		api.POST("/create", createUserCtrl.Handle)
	}
}