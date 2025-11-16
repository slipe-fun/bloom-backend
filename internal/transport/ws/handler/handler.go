package handler

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

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
				events.Send(hub, client, clientToken, session.UserID, room, socketMsg)
			case "add_keys":
				var socketKeys domain.SocketKeys
				if err := json.Unmarshal(msg, &socketKeys); err != nil {
					events.SendError(client, "invalid_message_format")
					continue
				}
				events.AddChatKeys(hub, client, clientToken, session.UserID, socketKeys)
			case "create_chat":
				var socketChat domain.SocketChat
				if err := json.Unmarshal(msg, &socketChat); err != nil {
					events.SendError(client, "invalid_message_format")
					continue
				}

				events.CreateChat(hub, client, clientToken, session.UserID, socketChat)

				chats, _ = hub.Chats.GetChatsByUserId(clientToken)
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
				events.MessageSeen(hub, client, clientToken, session.UserID, room, seenMsg)

			default:
				log.Println("Unknown message type:", baseMsg.Type)
			}
		}
	}
}
