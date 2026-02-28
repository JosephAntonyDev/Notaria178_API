package controllers

import (
	"net/http"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

type RemoveCollaboratorController struct {
	useCase *app.RemoveCollaboratorUseCase
}

func NewRemoveCollaboratorController(uc *app.RemoveCollaboratorUseCase) *RemoveCollaboratorController {
	return &RemoveCollaboratorController{useCase: uc}
}

func (ctrl *RemoveCollaboratorController) Handle(c *gin.Context) {
	reqCtx := extractRequestContext(c)
	workID := c.Param("id")
	userID := c.Param("userId")

	if err := ctrl.useCase.Execute(c.Request.Context(), reqCtx, workID, userID); err != nil {
		handleUseCaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Colaborador removido exitosamente"})
}
