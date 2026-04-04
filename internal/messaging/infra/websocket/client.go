package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// Client representa una conexión WebSocket individual
type Client struct {
	ID    string
	Conn  *websocket.Conn
	Send  chan []byte
	Rooms map[string]bool
	Hub   *Hub
	mu    sync.Mutex
}

// NewClient crea un nuevo cliente WebSocket
func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:    id,
		Conn:  conn,
		Send:  make(chan []byte, 256),
		Rooms: make(map[string]bool),
		Hub:   hub,
	}
}

// MessageHandler define la interface para manejar mensajes entrantes
type MessageHandler interface {
	HandleMessage(client *Client, message []byte)
}

// ReadPump lee mensajes del WebSocket
func (c *Client) ReadPump(handler MessageHandler) {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		handler.HandleMessage(c, message)
	}
}

// WritePump escribe mensajes al WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage envía un mensaje serializado al cliente
func (c *Client) SendMessage(data []byte) {
	select {
	case c.Send <- data:
	default:
		log.Printf("Failed to send message to client %s: buffer full", c.ID)
	}
}
