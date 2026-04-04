package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Tipos de mensajes WebSocket
const (
	WSTypeComment   = "comment"
	WSTypeTyping    = "typing"
	WSTypeJoinRoom  = "join_room"
	WSTypeLeaveRoom = "leave_room"
	WSTypeError     = "error"
	WSTypeHistory   = "history"
)

// WSMessage es el payload que viaja por WebSocket
type WSMessage struct {
	Type      string      `json:"type"`
	RoomID    string      `json:"room_id,omitempty"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

// CommentPayload es el payload específico para comentarios
type CommentPayload struct {
	ID        uuid.UUID `json:"id"`
	WorkID    uuid.UUID `json:"work_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// IncomingMessage representa un mensaje entrante del cliente
type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// JoinRoomPayload es el payload para unirse a una sala
type JoinRoomPayload struct {
	WorkID string `json:"work_id"`
}

// SendCommentPayload es el payload para enviar un comentario
type SendCommentPayload struct {
	WorkID  string `json:"work_id"`
	Message string `json:"message"`
}
