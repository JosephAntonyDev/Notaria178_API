package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type RemoveWorkActController struct {
	useCase *app.RemoveWorkActUseCase
}

func NewRemoveWorkActController(uc *app.RemoveWorkActUseCase) *RemoveWorkActController {
	return &RemoveWorkActController{useCase: uc}
}

func (ctrl *RemoveWorkActController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")
	actID := c.Param("actId")

	err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, actID)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Acto desasociado exitosamente",
	})
}
