package app

import (
	"context"
	"time"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/repository"
	"github.com/google/uuid"
)

type RegisterDeviceTokenRequest struct {
	FCMToken   string `json:"fcm_token" binding:"required"`
	DeviceType string `json:"device_type"` // "web", "android", "ios"
}

type RegisterDeviceTokenUseCase struct {
	repo repository.DeviceTokenRepository
}

func NewRegisterDeviceTokenUseCase(repo repository.DeviceTokenRepository) *RegisterDeviceTokenUseCase {
	return &RegisterDeviceTokenUseCase{
		repo: repo,
	}
}

func (uc *RegisterDeviceTokenUseCase) Execute(ctx context.Context, userID uuid.UUID, req RegisterDeviceTokenRequest) error {
	// Determinar el tipo de dispositivo
	deviceType := entities.DeviceTypeWeb
	if req.DeviceType == "android" {
		deviceType = entities.DeviceTypeAndroid
	} else if req.DeviceType == "ios" {
		deviceType = entities.DeviceTypeIOS
	}

	token := &entities.DeviceToken{
		ID:         uuid.New(),
		UserID:     userID,
		FCMToken:   req.FCMToken,
		DeviceType: deviceType,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// El repositorio hace upsert, así que si el token ya existe, solo actualiza
	if err := uc.repo.SaveToken(ctx, token); err != nil {
		return err
	}

	return nil
}
