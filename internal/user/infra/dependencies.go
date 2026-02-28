package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/adapters"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/infra/routers"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	userRepo := repository.NewPostgresUserRepository(db)
	hasher := adapters.NewBcrypt()
	tokenManager := adapters.NewJWTManager(jwtSecret)

	createUserUseCase := app.NewCreateUserUseCase(userRepo, hasher)
	loginUserUseCase := app.NewLoginUserUseCase(userRepo, hasher, tokenManager)

	createUserCtrl := controllers.NewCreateUserController(createUserUseCase)
	loginUserCtrl := controllers.NewLoginUserController(loginUserUseCase)

	routers.SetupUserRoutes(r, createUserCtrl, loginUserCtrl)
}