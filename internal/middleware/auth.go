package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

// 1. AuthMiddleware: Valida que el token sea real y extrae los datos
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Falta el token de autorización"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma fraudulento")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Estructura del token corrupta"})
			return
		}

		// Extraemos y guardamos en el Contexto de Gin (usamos strings porque son UUIDs y Enums)
		if sub, ok := claims["sub"].(string); ok { // sub = Subject (ID del usuario)
			c.Set("userID", sub)
		}
		
		if role, ok := claims["role"].(string); ok {
			c.Set("userRole", role)
		}
		
		if branchID, ok := claims["branch_id"].(string); ok {
			c.Set("branchID", branchID)
		}

		c.Next()
	}
}

// 2. RoleMiddleware: Un middleware dinámico que bloquea a los intrusos
// Le pasas una lista de los roles permitidos para la ruta.
func RequireRoles(allowedRoles ...entities.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtenemos el rol que el AuthMiddleware guardó
		userRoleStr, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No se pudo determinar tu rol"})
			return
		}

		userRole := entities.UserRole(userRoleStr.(string))
		isAllowed := false

		// Revisamos si su rol está en la lista de invitados VIP
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "No tienes nivel de autorización para realizar esta acción",
			})
			return
		}

		c.Next()
	}
}