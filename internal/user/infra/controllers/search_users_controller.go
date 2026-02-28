package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/JosephAntonyDev/Notaria178_API/internal/core/dtos"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/app"
)

type SearchUsersController struct {
	useCase *app.SearchUsersUseCase
}

func NewSearchUsersController(uc *app.SearchUsersUseCase) *SearchUsersController {
	return &SearchUsersController{
		useCase: uc,
	}
}

func (ctrl *SearchUsersController) Handle(c *gin.Context) {
	//Instanciamos el DTO Global
	var req dtos.PaginationRequest

	// Valores por defecto seguros por si el frontend no los envía
	req.Limit = 10
	req.Offset = 0

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de paginación inválidos"})
		return
	}

	users, err := ctrl.useCase.Execute(c.Request.Context(), req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno al buscar los usuarios"})
		return
	}

	// Nota: Por ahora 'Total' es la longitud del arreglo actual. 
	// Para un conteo real de todos los registros (ej. "página 1 de 50"), 
	// en el futuro podemos agregar un método Count() al repositorio.
	response := dtos.PaginatedResponse{
		Total: len(users), 
		Data:  users,
	}

	c.JSON(http.StatusOK, response)
}