package events

import (
	"encoding/json"
	"log"

	"github.com/fasthttp/websocket"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func Send(hub *types.Hub, sender *types.Client, senderID int, room string, message domain.SocketMessage) {
	if clients, ok := hub.Clients[room]; ok {
		for client := range clients {
			outMsg := struct {
				UserID int `json:"user_id"`
				domain.SocketMessage
			}{
				UserID:        senderID,
				SocketMessage: message,
			}

			b, err := json.Marshal(outMsg)
			if err != nil {
				log.Println("Failed to marshal message:", err)
				continue
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Println("Failed to send message:", err)
			}
		}
	}
}
