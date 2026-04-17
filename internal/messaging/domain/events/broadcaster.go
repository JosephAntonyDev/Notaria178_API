package events

import "github.com/JosephAntonyDev/Notaria178_API/internal/messaging/domain/entities"

// Broadcaster es el puerto de salida para emitir mensajes a salas
type Broadcaster interface {
	BroadcastToRoom(roomID string, message *entities.WSMessage)
	AddClientToRoom(clientID string, roomID string)
	RemoveClientFromRoom(clientID string, roomID string)
}
