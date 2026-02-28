package controllers

import (
	"io"

	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/notification/infra/events"
	"github.com/gin-gonic/gin"
)

// StreamNotificationsController mantiene una conexión SSE abierta por usuario.
type StreamNotificationsController struct {
	hub *events.SSEHub
}

func NewStreamNotificationsController(hub *events.SSEHub) *StreamNotificationsController {
	return &StreamNotificationsController{hub: hub}
}

func (ctrl *StreamNotificationsController) Handle(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(string)

	// Registrar al cliente en el hub
	clientChan := ctrl.hub.AddClient(uid)
	defer ctrl.hub.RemoveClient(uid)

	// Headers obligatorios para SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Mantener la conexión abierta, escuchando eventos
	c.Stream(func(w io.Writer) bool {
		select {
		case notif, ok := <-clientChan:
			if !ok {
				// Canal cerrado, terminar stream
				return false
			}
			dto := app.ToNotificationDTO(notif)
			c.SSEvent("notification", dto)
			return true
		case <-c.Request.Context().Done():
			// El cliente cerró la conexión
			return false
		}
	})
}
