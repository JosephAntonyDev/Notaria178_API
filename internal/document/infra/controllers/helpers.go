package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// handleUseCaseError mapea los errores de la capa de aplicación a códigos HTTP
func handleUseCaseError(c *gin.Context, err error) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "inválido") || strings.Contains(msg, "inválida"):
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	case strings.Contains(msg, "no encontrado"):
		c.JSON(http.StatusNotFound, gin.H{"error": msg})
	case strings.Contains(msg, "no tienes") || strings.Contains(msg, "no puedes"):
		c.JSON(http.StatusForbidden, gin.H{"error": msg})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
	}
}
