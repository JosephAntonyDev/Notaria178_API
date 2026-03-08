package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type AddCommentController struct {
	useCase *app.AddCommentUseCase
}

func NewAddCommentController(uc *app.AddCommentUseCase) *AddCommentController {
	return &AddCommentController{useCase: uc}
}

func (ctrl *AddCommentController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	var req app.AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: el mensaje es requerido"})
		return
	}

	comment, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, req)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comentario agregado exitosamente",
		"data":    comment,
	})
}
