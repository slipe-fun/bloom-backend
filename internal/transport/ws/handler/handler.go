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
				for room := range client.Rooms {
					events.Leave(hub, client, room)
				}
				break
			}

			var baseMsg struct {
				Type   string `json:"type"`
				ChatID int    `json:"chat_id"`
			}
			if err := json.Unmarshal(msg, &baseMsg); err != nil {
				events.SendError(client, "invalid_message_format")
				continue
			}

			switch baseMsg.Type {
			case "send":
				var socketMsg domain.SocketMessage
				if err := json.Unmarshal(msg, &socketMsg); err != nil {
					events.SendError(client, "invalid_message_format")
					continue
				}

				room := "chat" + strconv.Itoa(socketMsg.ChatID)
				if !service.IsUserInChat(chats, socketMsg.ChatID) {
					events.SendError(client, "not_member")
					continue
				}
				events.Send(hub, client, clientToken, userID, room, socketMsg)

			case "message_seen":
				var seenMsg domain.SocketMessageSeen
				if err := json.Unmarshal(msg, &seenMsg); err != nil {
					events.SendError(client, "invalid_message_format")
					continue
				}

				room := "chat" + strconv.Itoa(seenMsg.ChatID)
				if !service.IsUserInChat(chats, seenMsg.ChatID) {
					events.SendError(client, "not_member")
					continue
				}
				events.MessageSeen(hub, client, clientToken, userID, room, seenMsg)

			default:
				log.Println("Unknown message type:", baseMsg.Type)
			}
		}
	}
}
