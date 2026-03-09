package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/app"
	"github.com/gin-gonic/gin"
)

type DeleteDocumentController struct {
	useCase *app.DeleteDocumentUseCase
}

func NewDeleteDocumentController(uc *app.DeleteDocumentUseCase) *DeleteDocumentController {
	return &DeleteDocumentController{useCase: uc}
}

func (ctrl *DeleteDocumentController) Handle(c *gin.Context) {
	docID := c.Param("id")

	if err := ctrl.useCase.Execute(c.Request.Context(), docID); err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Documento eliminado exitosamente",
	})
}
