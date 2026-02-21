package types

import (
	"github.com/gofiber/websocket/v2"
	ChatApp "github.com/slipe-fun/skid-backend/internal/app/chat"
	SessionApp "github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

type Client struct {
	Conn   *websocket.Conn
	Rooms  map[string]bool
	UserID int
}

type Hub struct {
	Clients         map[string]map[*Client]bool
	ClientsByUserID map[int]map[*Client]bool
	SessionApp      *SessionApp.SessionApp
	Chats           *ChatApp.ChatApp
}

func NewHub(sessionApp *SessionApp.SessionApp, chats *ChatApp.ChatApp) *Hub {
	return &Hub{
		SessionApp:      sessionApp,
		Clients:         make(map[string]map[*Client]bool),
		ClientsByUserID: make(map[int]map[*Client]bool),
		Chats:           chats,
	}
}

func (h *Hub) JoinRoom(client *Client, room string) {
	if client.Rooms == nil {
		client.Rooms = make(map[string]bool)
	}
	client.Rooms[room] = true

	if _, ok := h.Clients[room]; !ok {
		h.Clients[room] = make(map[*Client]bool)
	}
	h.Clients[room][client] = true
}

func (h *Hub) LeaveRoom(client *Client, room string) {
	delete(client.Rooms, room)

	if clients, ok := h.Clients[room]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.Clients, room)
		}
	}
}

func (h *Hub) LeaveAllRooms(client *Client) {
	for room := range client.Rooms {
		h.LeaveRoom(client, room)
	}
}

func (h *Hub) RegisterUser(userID int, client *Client) {
	if h.ClientsByUserID[userID] == nil {
		h.ClientsByUserID[userID] = make(map[*Client]bool)
	}
	h.ClientsByUserID[userID][client] = true
	metrics.ActiveWebsocketConnections.Inc()
}

func (h *Hub) UnregisterUser(userID int, client *Client) {
	if clients, ok := h.ClientsByUserID[userID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.ClientsByUserID, userID)
		}
	}
	metrics.ActiveWebsocketConnections.Dec()
}

func (h *Hub) SendToUser(userID int, message []byte) {
	if clients, ok := h.ClientsByUserID[userID]; ok {
		for client := range clients {
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				h.UnregisterUser(userID, client)
				client.Conn.Close()
			}
		}
	}
}

func (h *Hub) Broadcast(room string, message []byte) {
	if clients, ok := h.Clients[room]; ok {
		for client := range clients {
			client.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
