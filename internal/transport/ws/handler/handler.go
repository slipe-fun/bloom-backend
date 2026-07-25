package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/events"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

const (
	pongWait = 60 * time.Second
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

		client := types.NewClient(hub, c, session.UserID)

		hub.RegisterUser(session.UserID, client)
		go client.WritePump()

		defer func() {
			hub.UnregisterUser(session.UserID, client)
			hub.LeaveAllRooms(client)
			client.Close()
		}()

		chats, err := hub.Chats.GetChatsByUserID(session.UserID)
		if err == nil && len(chats) > 0 {
			for _, chat := range chats {
				events.Join(hub, client, "chat"+strconv.Itoa(chat.ID))
			}
		}

		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}
