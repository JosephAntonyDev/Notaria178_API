package controllers

import (
	"net/http"
	"os"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/app"
	"github.com/gin-gonic/gin"
)

type DownloadDocumentController struct {
	useCase *app.GetDocumentUseCase
}

func NewDownloadDocumentController(uc *app.GetDocumentUseCase) *DownloadDocumentController {
	return &DownloadDocumentController{useCase: uc}
}

func (ctrl *DownloadDocumentController) Handle(c *gin.Context) {
	docID := c.Param("id")

	doc, err := ctrl.useCase.Execute(c.Request.Context(), docID)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	// Verificar que el archivo físico exista en disco
	if _, err := os.Stat(doc.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "El archivo físico no fue encontrado en el servidor"})
		return
	}

	c.File(doc.FilePath)
}
