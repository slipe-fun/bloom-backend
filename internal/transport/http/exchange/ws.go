// internal/transport/http/exchange/ws.go
package exchange

import (
	"context"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/events"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func HandleExchangeWS(hub *types.Hub, rdb *redis.Client) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()

		client := &types.Client{
			Conn:   c,
			UserID: 0,
		}

		roomID := c.Query("room_id")
		if roomID == "" {
			events.SendError(client, "missing room id")
			return
		}

		ctx := context.Background()

		remaining, err := rdb.Decr(ctx, "exchange:session:"+roomID).Result()
		if err != nil {
			events.SendError(client, "failed to validate room")
			return
		}

		if remaining < 0 {
			events.SendError(client, "room is full, invalid or expired")
			return
		}

		if remaining == 0 {
			_ = rdb.Del(ctx, "exchange:session:"+roomID)
		}

		roomName := "exchange:" + roomID
		hub.JoinRoom(client, roomName)

		defer func() {
			hub.LeaveRoom(client, roomName)
			c.Close()
		}()

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
