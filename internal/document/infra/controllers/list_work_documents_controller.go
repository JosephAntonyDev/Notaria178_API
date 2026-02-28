package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/document/app"
	"github.com/gin-gonic/gin"
)

type ListWorkDocumentsController struct {
	useCase *app.ListWorkDocumentsUseCase
}

func NewListWorkDocumentsController(uc *app.ListWorkDocumentsUseCase) *ListWorkDocumentsController {
	return &ListWorkDocumentsController{useCase: uc}
}

func (ctrl *ListWorkDocumentsController) Handle(c *gin.Context) {
	workID := c.Param("work_id")

	docs, err := ctrl.useCase.Execute(c.Request.Context(), workID)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": docs})
}
