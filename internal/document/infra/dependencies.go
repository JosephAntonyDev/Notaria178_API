package infra

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/repository"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/routes"
	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/storage"
)

func SetupDependencies(r *gin.Engine, db *sql.DB, jwtSecret string) {
	docRepo := repository.NewPostgresDocumentRepository(db)
	fileStorage := storage.NewLocalFileStorage()

	// Casos de uso
	uploadDocUC := app.NewUploadDocumentUseCase(docRepo, fileStorage)
	listWorkDocsUC := app.NewListWorkDocumentsUseCase(docRepo)
	getDocUC := app.NewGetDocumentUseCase(docRepo)

	// Controladores
	uploadCtrl := controllers.NewUploadDocumentController(uploadDocUC)
	listWorkDocsCtrl := controllers.NewListWorkDocumentsController(listWorkDocsUC)
	downloadCtrl := controllers.NewDownloadDocumentController(getDocUC)

	routes.SetupDocumentRoutes(r, uploadCtrl, listWorkDocsCtrl, downloadCtrl, jwtSecret)
}
