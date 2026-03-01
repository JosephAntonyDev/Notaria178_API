package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type GetCommentsController struct {
	useCase *app.ListCommentsUseCase
}

func NewGetCommentsController(uc *app.ListCommentsUseCase) *GetCommentsController {
	return &GetCommentsController{useCase: uc}
}

func (ctrl *GetCommentsController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	comments, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}
