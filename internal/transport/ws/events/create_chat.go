package events

import (
	"encoding/json"

	"github.com/fasthttp/websocket"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/transport/ws/types"
)

func CreateChat(hub *types.Hub, sender *types.Client, token string, senderId int, chat domain.SocketChat) {
	chatExists, err := hub.Chats.GetChatWithUsers(token, chat.Recipient)
	if err == nil && chatExists != nil {
		SendError(sender, "chat_already_exists")
		return
	}

	recipient, err := hub.Users.GetUserById(chat.Recipient)
	if err != nil {
		SendError(sender, "user_not_found")
		return
	}

	createChat, err := hub.Chats.CreateChat(token, recipient.ID)
	if err != nil {
		SendError(sender, "cant_create_chat")
		return
	}

	outMsg := struct {
		Type   string      `json:"type"`
		UserID int         `json:"user_id"`
		Chat   domain.Chat `json:"chat"`
	}{
		Type:   "chat_created",
		UserID: senderId,
		Chat:   *createChat,
	}

	b, err := json.Marshal(outMsg)
	if err != nil {
		return
	}
	room := "chat" + strconv.Itoa(createChat.ID)

	Join(hub, sender, room)

	if recipientClient, ok := hub.ClientsByUserID[chat.Recipient]; ok {
		Join(hub, recipientClient, room)
		recipientClient.Conn.WriteMessage(websocket.TextMessage, b)
	}

	sender.Conn.WriteMessage(websocket.TextMessage, b)
}
