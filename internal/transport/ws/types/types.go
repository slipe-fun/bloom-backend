package types

import (
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	ChatApp "github.com/slipe-fun/skid-backend/internal/app/chat"
	SessionApp "github.com/slipe-fun/skid-backend/internal/app/session"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

const (
	writeWait            = 10 * time.Second
	pingPeriod           = 30 * time.Second
	clientSendBufferSize = 256
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Rooms  map[string]bool
	UserID int
	Send   chan []byte
	Done   chan struct{}
	once   sync.Once
}

func NewClient(hub *Hub, conn *websocket.Conn, userID int) *Client {
	return &Client{
		Hub:    hub,
		Conn:   conn,
		Rooms:  make(map[string]bool),
		UserID: userID,
		Send:   make(chan []byte, clientSendBufferSize),
		Done:   make(chan struct{}),
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		close(c.Done)
		c.Conn.Close()
	})
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
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

		case <-c.Done:
			return
		}
	}
}

type Hub struct {
	mu              sync.RWMutex
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
	h.mu.Lock()
	defer h.mu.Unlock()

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
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(client.Rooms, room)

	if clients, ok := h.Clients[room]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.Clients, room)
		}
	}
}

func (h *Hub) LeaveAllRooms(client *Client) {
	h.mu.Lock()
	rooms := make([]string, 0, len(client.Rooms))
	for room := range client.Rooms {
		rooms = append(rooms, room)
	}
	h.mu.Unlock()

	for _, room := range rooms {
		h.LeaveRoom(client, room)
	}
}

func (h *Hub) RegisterUser(userID int, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.ClientsByUserID[userID] == nil {
		h.ClientsByUserID[userID] = make(map[*Client]bool)
	}
	h.ClientsByUserID[userID][client] = true
	metrics.ActiveWebsocketConnections.Inc()
}

func (h *Hub) UnregisterUser(userID int, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.ClientsByUserID[userID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.ClientsByUserID, userID)
		}
	}
	metrics.ActiveWebsocketConnections.Dec()
}

func (h *Hub) SendToUser(userID int, message []byte) {
	h.mu.RLock()
	clients, ok := h.ClientsByUserID[userID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	activeClients := make([]*Client, 0, len(clients))
	for client := range clients {
		activeClients = append(activeClients, client)
	}
	h.mu.RUnlock()

	for _, client := range activeClients {
		select {
		case client.Send <- message:
		default:
			client.Close()
		}
	}
}

func (h *Hub) Broadcast(room string, message []byte) {
	h.mu.RLock()
	clients, ok := h.Clients[room]
	if !ok {
		h.mu.RUnlock()
		return
	}

	activeClients := make([]*Client, 0, len(clients))
	for client := range clients {
		activeClients = append(activeClients, client)
	}
	h.mu.RUnlock()

	for _, client := range activeClients {
		select {
		case client.Send <- message:
		default:
			client.Close()
		}
	}
}
