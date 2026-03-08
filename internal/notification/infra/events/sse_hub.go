package events

import (
	"sync"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/domain/entities"
)

// SSEHub es la implementación de infraestructura del NotificationNotifier.
// Gestiona un mapa concurrente de canales, uno por usuario conectado.
type SSEHub struct {
	clients map[string]chan *entities.Notification
	mu      sync.RWMutex
}

// NewSSEHub crea una nueva instancia del hub de Server-Sent Events.
func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients: make(map[string]chan *entities.Notification),
	}
}

// AddClient registra un canal para el usuario y lo devuelve para que el controlador lo escuche.
func (h *SSEHub) AddClient(userID string) chan *entities.Notification {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Si ya existe un canal previo, cerrarlo antes de reemplazar
	if ch, ok := h.clients[userID]; ok {
		close(ch)
	}

	ch := make(chan *entities.Notification, 16)
	h.clients[userID] = ch
	return ch
}

// RemoveClient cierra y elimina el canal del usuario del mapa.
func (h *SSEHub) RemoveClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ch, ok := h.clients[userID]; ok {
		close(ch)
		delete(h.clients, userID)
	}
}

// Broadcast envía una notificación al canal del usuario si está conectado.
// Utiliza un envío no bloqueante para no detener al emisor si el canal está lleno.
func (h *SSEHub) Broadcast(userID string, notif *entities.Notification) {
	h.mu.RLock()
	ch, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	select {
	case ch <- notif:
	default:
		// Canal lleno; se descarta para no bloquear al emisor.
	}
}
