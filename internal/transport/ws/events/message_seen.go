package events

import (
	"encoding/json"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func MessageSeen(hub *types.Hub, sender *types.Client, token string, senderID int, room string, messageSeen domain.SocketMessageSeen) {
	if clients, ok := hub.Clients[room]; ok {
		if len(messageSeen.Messages) == 0 {
			return
		}

		var successMessages []int
		for i := range messageSeen.Messages {
			message, err := hub.Messages.GetMessageById(token, messageSeen.Messages[i])
			if err != nil {
				continue
			}

			if message.ChatID != messageSeen.ChatID || message.Seen != nil {
				continue
			}

			successMessages = append(successMessages, messageSeen.Messages[i])
		}

		seenLocalTime := time.Now()
		updateMessageError := hub.Messages.UpdateMessagesSeenStatus(successMessages, seenLocalTime)
		if updateMessageError != nil {
			SendError(sender, "failed_seen_message")
			return
		}

		outMsg := struct {
			UserID   int       `json:"user_id"`
			ChatID   int       `json:"chat_id"`
			SeenAt   time.Time `json:"seen_at"`
			Messages []int     `json:"messages"`
		}{
			UserID:   senderID,
			ChatID:   messageSeen.ChatID,
			SeenAt:   time.Now(),
			Messages: successMessages,
		}

		if len(outMsg.Messages) == 0 {
			SendError(sender, "failed_seen_message")
			return
		}

		b, err := json.Marshal(outMsg)
		if err != nil {
			return
		}

		for client := range clients {
			if err := client.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
				SendError(sender, "failed_seen_message")
			}
		}
	}
}
