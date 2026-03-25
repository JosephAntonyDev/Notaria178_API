package infra

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/adapters"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/events"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/routes"
)

// SetupResult agrupa los use cases que este modulo expone para integracion cruzada.
type SetupResult struct {
	CreateNotifUC      *app.CreateNotificationUseCase
	NotifyNewCommentUC *app.NotifyNewCommentUseCase
}

// SetupDependencies configura todo el modulo de notificaciones.
// workCollabGetter es la capacidad de consultar colaboradores de un trabajo.
// Se acepta como interfaz para evitar importar work/infra/repository directamente.
func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string, workCollabGetter adapters.CollaboratorGetter) *SetupResult {
	// Repositorios
	notifRepo := repository.NewPostgresNotificationRepository(db)
	deviceTokenRepo := repository.NewPostgresDeviceTokenRepository(db)

	// SSE Hub para notificaciones en tiempo real
	hub := events.NewSSEHub()

	// Firebase FCM Sender (opcional, solo si esta configurado)
	firebaseCredPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	var fcmSender *adapters.FirebaseNotificationSender
	if firebaseCredPath != "" {
		var err error
		fcmSender, err = adapters.NewFirebaseNotificationSender(firebaseCredPath)
		if err != nil {
			log.Printf("[WARN] No se pudo inicializar Firebase FCM: %v", err)
			log.Println("Las notificaciones push no estaran disponibles")
		}
	} else {
		log.Println("[WARN] FIREBASE_CREDENTIALS_PATH no configurado. Las notificaciones push no estaran disponibles")
	}

	// Casos de uso
	getMyNotificationsUC := app.NewGetMyNotificationsUseCase(notifRepo)
	markAsReadUC := app.NewMarkAsReadUseCase(notifRepo)
	markAllReadUC := app.NewMarkAllReadUseCase(notifRepo)
	createNotifUC := app.NewCreateNotificationUseCase(notifRepo, hub)
	registerDeviceTokenUC := app.NewRegisterDeviceTokenUseCase(deviceTokenRepo)
	getUnreadCountUC := app.NewGetUnreadCountUseCase(notifRepo)

	// Adaptador para que NotifyNewComment pueda consultar colaboradores
	workRepoAdapter := adapters.NewWorkRepoAdapter(workCollabGetter)

	// Use case de notificaciones de comentarios (push + in-app)
	var notifyNewCommentUC *app.NotifyNewCommentUseCase
	if fcmSender != nil {
		notifyNewCommentUC = app.NewNotifyNewCommentUseCase(
			notifRepo, deviceTokenRepo, workRepoAdapter, fcmSender, hub,
		)
	} else {
		// Sin FCM, crear un use case que solo haga in-app + SSE (sin push)
		notifyNewCommentUC = app.NewNotifyNewCommentUseCase(
			notifRepo, deviceTokenRepo, workRepoAdapter, &noopFCMSender{}, hub,
		)
	}

	// Controladores
	getMyNotificationsCtrl := controllers.NewGetMyNotificationsController(getMyNotificationsUC)
	markAsReadCtrl := controllers.NewMarkAsReadController(markAsReadUC)
	markAllReadCtrl := controllers.NewMarkAllReadController(markAllReadUC)
	sseCtrl := controllers.NewStreamNotificationsController(hub)
	registerDeviceTokenCtrl := controllers.NewRegisterDeviceTokenController(registerDeviceTokenUC)
	getUnreadCountCtrl := controllers.NewGetUnreadCountController(getUnreadCountUC)

	// Rutas
	routes.SetupNotificationRoutes(
		r,
		getMyNotificationsCtrl,
		markAsReadCtrl,
		markAllReadCtrl,
		sseCtrl,
		registerDeviceTokenCtrl,
		getUnreadCountCtrl,
		jwtSecret,
	)

	return &SetupResult{
		CreateNotifUC:      createNotifUC,
		NotifyNewCommentUC: notifyNewCommentUC,
	}
}
