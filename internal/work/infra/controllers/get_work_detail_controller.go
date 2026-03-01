package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type GetWorkDetailController struct {
	useCase *app.GetWorkDetailUseCase
}

func NewGetWorkDetailController(uc *app.GetWorkDetailUseCase) *GetWorkDetailController {
	return &GetWorkDetailController{useCase: uc}
}

func (ctrl *GetWorkDetailController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")

	work, err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": work})
}
