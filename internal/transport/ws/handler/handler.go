package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/events"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func HandleWS(hub *types.Hub) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()
		clientToken := c.Query("token")

		session, err := hub.SessionApp.GetSession(clientToken)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
			return
		}

		client := &types.Client{Conn: c}

		if hub.ClientsByUserID == nil {
			hub.ClientsByUserID = make(map[int]*types.Client)
		}
		hub.ClientsByUserID[session.UserID] = client

		hub.RegisterUser(session.UserID, client)

		defer func() {
			hub.UnregisterUser(session.UserID)
			hub.LeaveAllRooms(client)
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}

		chats, err := hub.Chats.GetChatsByUserId(clientToken)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Get chats error"))
			return
		}

		if len(chats) > 0 {
			for _, chat := range chats {
				events.Join(hub, client, "chat"+strconv.Itoa(chat.ID))
			}
		}

		c.SetReadDeadline(time.Now().Add(70 * time.Second))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(70 * time.Second))
			return nil
		})

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
					return
				}
			}
		}()
	}
}
