package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type DeleteWorkRequirementController struct {
	useCase *app.DeleteWorkRequirementUseCase
}

func NewDeleteWorkRequirementController(uc *app.DeleteWorkRequirementUseCase) *DeleteWorkRequirementController {
	return &DeleteWorkRequirementController{useCase: uc}
}

func (ctrl *DeleteWorkRequirementController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")
	reqID := c.Param("reqId")

	err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, reqID)
	if err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Requisito eliminado exitosamente",
	})
}
