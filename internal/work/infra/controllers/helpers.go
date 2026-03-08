package controllers

import (
	"net/http"
	"strings"

	"github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/gin-gonic/gin"
)

// extractRequestContext extrae userID, userRole y branchID del contexto de Gin (set por AuthMiddleware)
func extractRequestContext(c *gin.Context) app.RequestContext {
	reqCtx := app.RequestContext{}
	if uid, exists := c.Get("userID"); exists {
		reqCtx.UserID = uid.(string)
	}
	if role, exists := c.Get("userRole"); exists {
		reqCtx.UserRole = role.(string)
	}
	if bid, exists := c.Get("branchID"); exists {
		reqCtx.BranchID = bid.(string)
	}
	return reqCtx
}

// handleUseCaseError mapea los errores de la capa de aplicación a códigos HTTP
func handleUseCaseError(c *gin.Context, err error) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "inválido"):
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	case strings.Contains(msg, "no encontrado"):
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
	case strings.Contains(msg, "no tienes") || strings.Contains(msg, "no puedes"):
		c.JSON(http.StatusForbidden, gin.H{"error": msg})
	case strings.Contains(msg, "no permitida"):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
	}
}
