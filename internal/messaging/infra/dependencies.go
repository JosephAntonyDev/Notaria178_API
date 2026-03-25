package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/infra/routes"
	ws "github.com/JosephAntonyDev/Notaria178_API/internal/messaging/infra/websocket"
	workApp "github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	workEvents "github.com/JosephAntonyDev/Notaria178_API/internal/work/domain/events"
	workRepo "github.com/JosephAntonyDev/Notaria178_API/internal/work/infra/repository"
)

func SetupDependencies(
	r *gin.Engine,
	db *sql.DB,
	jwtSecret string,
	auditAdapter workEvents.AuditLogger,
) *ws.Hub {
	// Repositorio de work (necesario para validar acceso)
	workRepository := workRepo.NewPostgresWorkRepository(db)

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Use cases existentes de work
	addCommentUC := workApp.NewAddCommentUseCase(workRepository, auditAdapter)
	listCommentsUC := workApp.NewListCommentsUseCase(workRepository)

	// Use cases de messaging
	joinRoomUC := app.NewJoinRoomUseCase(workRepository)
	sendCommentUC := app.NewSendCommentUseCase(addCommentUC, hub)

	// Controller
	wsCtrl := controllers.NewWSController(hub, joinRoomUC, sendCommentUC, listCommentsUC)

	routes.SetupMessagingRoutes(r, wsCtrl, jwtSecret)

	return hub
}
