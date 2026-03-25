package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/domain/entities"
)

// RoomMessage representa un mensaje para broadcast a una sala
type RoomMessage struct {
	RoomID  string
	Message []byte
}

// Hub gestiona todas las conexiones WebSocket y salas
type Hub struct {
	clients    map[string]*Client         // userID -> Client
	rooms      map[string]map[string]bool // roomID -> set of userIDs
	register   chan *Client
	unregister chan *Client
	broadcast  chan *RoomMessage
	mu         sync.RWMutex
}

// NewHub crea una nueva instancia del hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		rooms:      make(map[string]map[string]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *RoomMessage, 256),
	}
}

// Run inicia el loop principal del hub (debe ejecutarse en una goroutine)
func (h *Hub) Run() {
	log.Println("WebSocket Hub started")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client %s connected", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				// Remover de todas las salas
				for roomID := range client.Rooms {
					if room, ok := h.rooms[roomID]; ok {
						delete(room, client.ID)
						if len(room) == 0 {
							delete(h.rooms, roomID)
						}
					}
				}
				close(client.Send)
				delete(h.clients, client.ID)
				log.Printf("Client %s disconnected", client.ID)
			}
			h.mu.Unlock()

		case rm := <-h.broadcast:
			h.mu.RLock()
			if room, ok := h.rooms[rm.RoomID]; ok {
				for userID := range room {
					if client, ok := h.clients[userID]; ok {
						select {
						case client.Send <- rm.Message:
						default:
							// Buffer lleno, cliente lento - cerrar conexión
							log.Printf("Client %s buffer full, closing", userID)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register expone el canal para registrar clientes
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister expone el canal para desregistrar clientes
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// BroadcastToRoom implementa la interface Broadcaster
func (h *Hub) BroadcastToRoom(roomID string, message *entities.WSMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	h.broadcast <- &RoomMessage{RoomID: roomID, Message: data}
}

// AddClientToRoom agrega un cliente a una sala
func (h *Hub) AddClientToRoom(clientID string, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[string]bool)
	}
	h.rooms[roomID][clientID] = true

	if client, ok := h.clients[clientID]; ok {
		client.mu.Lock()
		client.Rooms[roomID] = true
		client.mu.Unlock()
	}
	log.Printf("Client %s joined room %s", clientID, roomID)
}

// RemoveClientFromRoom remueve un cliente de una sala
func (h *Hub) RemoveClientFromRoom(clientID string, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		delete(room, clientID)
		if len(room) == 0 {
			delete(h.rooms, roomID)
		}
	}
	if client, ok := h.clients[clientID]; ok {
		client.mu.Lock()
		delete(client.Rooms, roomID)
		client.mu.Unlock()
	}
	log.Printf("Client %s left room %s", clientID, roomID)
}
