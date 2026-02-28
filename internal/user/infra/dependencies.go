package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/adapters"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	userRepo := repository.NewPostgresUserRepository(db)
	hasher := adapters.NewBcrypt()
	tokenManager := adapters.NewJWTManager(jwtSecret)

	createUserUseCase := app.NewCreateUserUseCase(userRepo, hasher)
	loginUserUseCase := app.NewLoginUserUseCase(userRepo, hasher, tokenManager)
	getProfileUseCase := app.NewGetProfileUseCase(userRepo)
	searchUsersUseCase := app.NewSearchUsersUseCase(userRepo)
	updateProdileUseCase := app.NewUpdateProfileUseCase(userRepo, hasher)
	updateEmployeeUseCase := app.NewUpdateEmployeeUseCase(userRepo, hasher)

	createUserCtrl := controllers.NewCreateUserController(createUserUseCase)
	loginUserCtrl := controllers.NewLoginUserController(loginUserUseCase)
	getProfileCtrl := controllers.NewGetProfileController(getProfileUseCase)
	searchUsersCtrl := controllers.NewSearchUsersController(searchUsersUseCase)
	updateProfileCtrl := controllers.NewUpdateProfileController(updateProdileUseCase)
	updateEmployeeCtrl := controllers.NewUpdateEmployeeController(updateEmployeeUseCase)

	routes.SetupUserRoutes(r, createUserCtrl, loginUserCtrl, getProfileCtrl, searchUsersCtrl, updateProfileCtrl, updateEmployeeCtrl, jwtSecret)
}