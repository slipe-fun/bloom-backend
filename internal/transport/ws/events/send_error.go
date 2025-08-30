package events

import (
	"encoding/json"

	"github.com/fasthttp/websocket"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func SendError(client *types.Client, errMsg string) {
	outErr := struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}{
		Type:    "message_error",
		Message: errMsg,
	}

	b, err := json.Marshal(outErr)
	if err != nil {
		return
	}

	_ = client.Conn.WriteMessage(websocket.TextMessage, b)
}
