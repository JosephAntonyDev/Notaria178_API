package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/infra/controllers"
	"github.com/JosephAntonyDev/Notaria178_API/internal/middleware"
)

func SetupDocumentRoutes(
	r *gin.Engine,
	uploadCtrl *controllers.UploadDocumentController,
	listWorkDocsCtrl *controllers.ListWorkDocumentsController,
	downloadCtrl *controllers.DownloadDocumentController,
	jwtSecret string,
) {
	api := r.Group("/api/v1/documents")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.POST("/upload", uploadCtrl.Handle)
		api.GET("/work/:work_id", listWorkDocsCtrl.Handle)
		api.GET("/download/:id", downloadCtrl.Handle)
	}
}
