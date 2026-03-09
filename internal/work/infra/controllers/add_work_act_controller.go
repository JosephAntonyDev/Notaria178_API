package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type AddWorkActController struct {
	useCase *app.AddWorkActUseCase
}

func NewAddWorkActController(uc *app.AddWorkActUseCase) *AddWorkActController {
	return &AddWorkActController{useCase: uc}
}

func (ctrl *AddWorkActController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	var req app.AddWorkActRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: se requiere act_id"})
		return
	}

	detail, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, req)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Acto asociado exitosamente",
		"data":    detail,
	})
}
