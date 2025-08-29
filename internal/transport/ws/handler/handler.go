package handler

import (
	"strings"

	"github.com/gofiber/websocket/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/events"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func HandleWS(hub *types.Hub) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()
		client := &types.Client{Conn: c}

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				events.Leave(hub, client)
				break
			}

			parts := strings.SplitN(string(msg), ":", 2)
			if len(parts) < 2 {
				continue
			}

			cmd := parts[0]
			data := parts[1]

			switch cmd {
			case "join":
				events.Join(hub, client, data)
			case "send":
				roomMsg := strings.SplitN(data, "|", 2)
				if len(roomMsg) == 2 {
					events.Send(hub, client, roomMsg[0], roomMsg[1])
				}
			case "leave":
				events.Leave(hub, client)
			}
		}
	}
}
