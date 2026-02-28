package adapters

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/JosephAntonyDev/Notaria178_API/internal/user/domain/entities"
)

type JWTManager struct {
	SecretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{SecretKey: secretKey}
}

func (j *JWTManager) GenerateToken(userID uuid.UUID, role entities.UserRole, branchID *uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": string(role),
		"iss":  "notaria178_api",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(12 * time.Hour).Unix(),
	}

	if branchID != nil {
		claims["branch_id"] = branchID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (bool, map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma fraudulento: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return false, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, claims, nil
	}

	return false, nil, fmt.Errorf("token inválido")
}