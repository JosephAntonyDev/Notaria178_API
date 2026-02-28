package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/routes"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	notifRepo := repository.NewPostgresNotificationRepository(db)
	hub := events.NewSSEHub()

	// Casos de uso
	getMyNotificationsUC := app.NewGetMyNotificationsUseCase(notifRepo)
	markAsReadUC := app.NewMarkAsReadUseCase(notifRepo)
	markAllReadUC := app.NewMarkAllReadUseCase(notifRepo)
	_ = app.NewCreateNotificationUseCase(notifRepo, hub)

	// Controladores
	getMyNotificationsCtrl := controllers.NewGetMyNotificationsController(getMyNotificationsUC)
	markAsReadCtrl := controllers.NewMarkAsReadController(markAsReadUC)
	markAllReadCtrl := controllers.NewMarkAllReadController(markAllReadUC)
	sseCtrl := controllers.NewStreamNotificationsController(hub)

	routes.SetupNotificationRoutes(r, getMyNotificationsCtrl, markAsReadCtrl, markAllReadCtrl, sseCtrl, jwtSecret)
}
