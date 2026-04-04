package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/app"
	"github.com/JosephAntonyDev/Notaria178_API/internal/messaging/domain/entities"
	ws "github.com/JosephAntonyDev/Notaria178_API/internal/messaging/infra/websocket"
	workApp "github.com/JosephAntonyDev/Notaria178_API/internal/work/app"
	"github.com/google/uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configurar en producción
	},
}

type WSController struct {
	hub            *ws.Hub
	joinRoomUC     *app.JoinRoomUseCase
	sendCommentUC  *app.SendCommentUseCase
	listCommentsUC *workApp.ListCommentsUseCase
}

func NewWSController(
	hub *ws.Hub,
	joinRoomUC *app.JoinRoomUseCase,
	sendCommentUC *app.SendCommentUseCase,
	listCommentsUC *workApp.ListCommentsUseCase,
) *WSController {
	return &WSController{
		hub:            hub,
		joinRoomUC:     joinRoomUC,
		sendCommentUC:  sendCommentUC,
		listCommentsUC: listCommentsUC,
	}
}

func (ctrl *WSController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	userRole, _ := c.Get("userRole")
	branchID, _ := c.Get("branchID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := ws.NewClient(userID.(string), conn, ctrl.hub)
	ctrl.hub.Register(client)

	reqCtx := workApp.RequestContext{
		UserID:   userID.(string),
		UserRole: userRole.(string),
	}
	if branchID != nil {
		reqCtx.BranchID = branchID.(string)
	}

	handler := &messageHandler{
		ctrl:   ctrl,
		client: client,
		reqCtx: reqCtx,
	}

	go client.WritePump()
	go client.ReadPump(handler)
}

type messageHandler struct {
	ctrl   *WSController
	client *ws.Client
	reqCtx workApp.RequestContext
}

func (h *messageHandler) HandleMessage(client *ws.Client, message []byte) {
	var incoming entities.IncomingMessage
	if err := json.Unmarshal(message, &incoming); err != nil {
		h.sendError("mensaje inválido")
		return
	}

	ctx := context.Background()

	switch incoming.Type {
	case entities.WSTypeJoinRoom:
		h.handleJoinRoom(ctx, incoming.Payload)

	case entities.WSTypeComment:
		h.handleComment(ctx, incoming.Payload)

	case entities.WSTypeLeaveRoom:
		h.handleLeaveRoom(incoming.Payload)
	}
}

func (h *messageHandler) handleJoinRoom(ctx context.Context, payload json.RawMessage) {
	var p entities.JoinRoomPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		h.sendError("payload inválido")
		return
	}

	workUUID, err := uuid.Parse(p.WorkID)
	if err != nil {
		h.sendError("ID de trabajo inválido")
		return
	}

	room := entities.NewWorkRoom(workUUID)

	if err := h.ctrl.joinRoomUC.Execute(ctx, h.reqCtx, room); err != nil {
		h.sendError(err.Error())
		return
	}

	h.ctrl.hub.AddClientToRoom(h.client.ID, room.ID)

	// Enviar historial de comentarios
	comments, err := h.ctrl.listCommentsUC.Execute(ctx, h.reqCtx, p.WorkID)
	if err != nil {
		h.sendError(err.Error())
		return
	}

	historyMsg := entities.WSMessage{
		Type:      entities.WSTypeHistory,
		RoomID:    room.ID,
		Payload:   comments,
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(historyMsg)
	h.client.SendMessage(data)
}

func (h *messageHandler) handleComment(ctx context.Context, payload json.RawMessage) {
	var p entities.SendCommentPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		h.sendError("payload inválido")
		return
	}

	if p.Message == "" {
		h.sendError("el mensaje no puede estar vacío")
		return
	}

	input := app.SendCommentInput{
		WorkID:  p.WorkID,
		Message: p.Message,
	}

	if err := h.ctrl.sendCommentUC.Execute(ctx, h.reqCtx, input); err != nil {
		h.sendError(err.Error())
	}
}

func (h *messageHandler) handleLeaveRoom(payload json.RawMessage) {
	var p entities.JoinRoomPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return
	}

	workUUID, err := uuid.Parse(p.WorkID)
	if err != nil {
		return
	}

	room := entities.NewWorkRoom(workUUID)
	h.ctrl.hub.RemoveClientFromRoom(h.client.ID, room.ID)
}

func (h *messageHandler) sendError(msg string) {
	errMsg := entities.WSMessage{
		Type:      entities.WSTypeError,
		Payload:   map[string]string{"error": msg},
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(errMsg)
	h.client.SendMessage(data)
}
