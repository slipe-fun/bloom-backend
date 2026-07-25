package exchange

import (
	"context"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

const (
	pongWait = 60 * time.Second
)

func HandleExchangeWS(hub *types.Hub, rdb *redis.Client) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()

		roomID := c.Query("room_id")
		if roomID == "" {
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"type":"message_error","message":"missing room id"}`))
			return
		}

		ctx := context.Background()

		remaining, err := rdb.Decr(ctx, "exchange:session:"+roomID).Result()
		if err != nil {
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"type":"message_error","message":"failed to validate room"}`))
			return
		}

		if remaining < 0 {
			_ = c.WriteMessage(websocket.TextMessage, []byte(`{"type":"message_error","message":"room is full, invalid or expired"}`))
			return
		}

		if remaining == 0 {
			_ = rdb.Del(ctx, "exchange:session:"+roomID)
		}

		client := types.NewClient(hub, c, 0)
		go client.WritePump()

		roomName := "exchange:" + roomID
		hub.JoinRoom(client, roomName)

		defer func() {
			hub.LeaveRoom(client, roomName)
			client.Close()
		}()

		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		for {
			messageType, msgBytes, err := c.ReadMessage()
			if err != nil {
				break
			}

			if messageType == websocket.TextMessage {
				hub.Broadcast(roomName, msgBytes)
			}
		}
	}
}
