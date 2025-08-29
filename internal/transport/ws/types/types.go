package types

import "github.com/gofiber/websocket/v2"

type Client struct {
	Conn *websocket.Conn
	Room string
}

type Hub struct {
	Clients map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{Clients: make(map[string]map[*Client]bool)}
}

func (h *Hub) JoinRoom(client *Client, room string) {
	client.Room = room
	if _, ok := h.Clients[room]; !ok {
		h.Clients[room] = make(map[*Client]bool)
	}
	h.Clients[room][client] = true
}

func (h *Hub) LeaveRoom(client *Client) {
	if clients, ok := h.Clients[client.Room]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.Clients, client.Room)
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
