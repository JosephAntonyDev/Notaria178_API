package entities

import "github.com/google/uuid"

type RoomType string

const (
	RoomTypeWorkComments RoomType = "WORK_COMMENTS"
	RoomTypePrivateChat  RoomType = "PRIVATE_CHAT" // Para uso futuro
)

// Room representa una sala de mensajes (puede ser un work o un chat privado)
type Room struct {
	ID       string // Para works: "work:{workID}", Para chat: "chat:{userA}:{userB}"
	Type     RoomType
	EntityID uuid.UUID // ID del work o ID del chat
}

func NewWorkRoom(workID uuid.UUID) *Room {
	return &Room{
		ID:       "work:" + workID.String(),
		Type:     RoomTypeWorkComments,
		EntityID: workID,
	}
}
