package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/app"
	"github.com/gin-gonic/gin"
)

type UploadDocumentController struct {
	useCase *app.UploadDocumentUseCase
}

func NewUploadDocumentController(uc *app.UploadDocumentUseCase) *UploadDocumentController {
	return &UploadDocumentController{useCase: uc}
}

func (ctrl *UploadDocumentController) Handle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo es requerido"})
		return
	}

	workID := c.PostForm("work_id")
	if workID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El work_id es requerido"})
		return
	}

	category := c.PostForm("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La categoría es requerida"})
		return
	}

	documentName := c.PostForm("document_name")
	if documentName == "" {
		documentName = file.Filename
	}

	// Extraer datos del JWT
	userID, _ := c.Get("userID")
	branchID, _ := c.Get("branchID")

	clientIDStr := c.PostForm("client_id")
	var clientID *string
	if clientIDStr != "" {
		clientID = &clientIDStr
	}

	input := app.UploadDocumentInput{
		File:         file,
		BranchID:     branchID.(string),
		WorkID:       workID,
		ClientID:     clientID,
		UserID:       userID.(string),
		DocumentName: documentName,
		Category:     category,
	}

	doc, err := ctrl.useCase.Execute(c.Request.Context(), input)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Documento subido exitosamente",
		"data":    doc,
	})
}
