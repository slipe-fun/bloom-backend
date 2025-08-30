package events

import (
	"encoding/json"

	"github.com/fasthttp/websocket"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func AddChatKeys(hub *types.Hub, sender *types.Client, token string, senderId int, keys domain.SocketKeys) {
	chat, err := hub.Chats.GetChatById(token, keys.ChatID)
	if err != nil {
		SendError(sender, "chat_not_found")
		return
	}

	if !hub.Chats.HasMember(chat, senderId) {
		SendError(sender, "not_member")
		return
	}

	keysCheckError := service.CheckKeysLength(keys.KyberPublicKey, keys.EcdhPublicKey, keys.EdPublicKey)
	if keysCheckError != nil {
		SendError(sender, keysCheckError.Error())
		return
	}

	updateChatError := hub.Chats.AddKeys(token, chat, keys.KyberPublicKey, keys.EcdhPublicKey, keys.EdPublicKey)
	if updateChatError != nil {
		SendError(sender, "cant_update_chat")
		return
	}

	recipient := 0
	for i := range chat.Members {
		if chat.Members[i].ID != senderId {
			recipient = chat.Members[i].ID
		}
	}

	outMsg := struct {
		Type   string `json:"type"`
		UserID int    `json:"user_id"`
		ChatID int    `json:"chat_id"`
		domain.SocketKeys
	}{
		Type:       "keys_added",
		UserID:     senderId,
		ChatID:     chat.ID,
		SocketKeys: keys,
	}

	b, err := json.Marshal(outMsg)
	if err != nil {
		return
	}

	if recipientClient, ok := hub.ClientsByUserID[recipient]; ok {
		recipientClient.Conn.WriteMessage(websocket.TextMessage, b)
	}
}
