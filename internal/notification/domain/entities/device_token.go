package entities

import (
	"time"

	"github.com/google/uuid"
)

// ─── Enum de tipo de dispositivo ────────────────────────────────────────────

type DeviceType string

const (
	DeviceTypeWeb     DeviceType = "web"
	DeviceTypeAndroid DeviceType = "android"
	DeviceTypeIOS     DeviceType = "ios"
)

// ─── Entidad DeviceToken (tabla: user_device_tokens) ────────────────────────

type DeviceToken struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	FCMToken   string     `json:"fcm_token"`
	DeviceType DeviceType `json:"device_type"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
