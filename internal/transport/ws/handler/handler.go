package handler

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/websocket/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/events"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func HandleWS(hub *types.Hub) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()
		clientToken := c.Query("token")
		_, err := hub.JwtSvc.VerifyToken(clientToken)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
			c.Close()
			return
		}

		userID, err := hub.TokenSvc.ExtractUserID(clientToken)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
			c.Close()
			return
		}

		client := &types.Client{Conn: c}

		chats, err := hub.Chats.GetChatsByUserId(clientToken)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Get chats error"))
			c.Close()
			return
		}

		if len(chats) > 0 {
			for _, chat := range chats {
				events.Join(hub, client, "chat"+strconv.Itoa(chat.ID))
			}
		}

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				events.Leave(hub, client)
				break
			}

			var socketMsg domain.SocketMessage
			if err := json.Unmarshal(msg, &socketMsg); err != nil {
				c.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
				continue
			}

			switch socketMsg.Type {
			case "send":
				room := "chat" + strconv.Itoa(socketMsg.ChatID)

				allowed := service.IsUserInChat(chats, socketMsg.ChatID)

				if !allowed {
					c.WriteMessage(websocket.TextMessage, []byte("You are not a member of this chat"))
					continue
				}

				events.Send(hub, client, userID, room, socketMsg)
			case "join":
				room := "chat" + strconv.Itoa(socketMsg.ChatID)

				allowed := service.IsUserInChat(chats, socketMsg.ChatID)

				if !allowed {
					c.WriteMessage(websocket.TextMessage, []byte("You are not a member of this chat"))
					continue
				}

				events.Join(hub, client, room)
			case "leave":
				events.Leave(hub, client)
			default:
				log.Println("Unknown message type:", socketMsg.Type)
			}
		}
	}
}
