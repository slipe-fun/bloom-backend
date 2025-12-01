package MessageApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *MessageApp) GetChatMessagesAfter(tokenStr string, chatId int, afterId int) ([]*domain.MessageWithReply, error) {
	_, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	chat, err := c.chats.GetChatById(tokenStr, chatId)
	if err != nil {
		return nil, err
	}

	messages, err := c.messages.GetChatMessagesAfter(chat.ID, afterId)
	if err != nil {
		return nil, domain.NotFound("messages not found")
	}

	result := make([]*domain.MessageWithReply, len(messages))
	for i, msg := range messages {
		messageWithReply := &domain.MessageWithReply{
			Message: *msg,
		}

		if msg.ReplyTo != nil {
			replyToMessage, err := c.messages.GetById(*msg.ReplyTo)
			if err == nil && replyToMessage != nil && replyToMessage.ChatID == msg.ChatID {
				messageWithReply.ReplyToMessage = replyToMessage
			}
		}

		result[i] = messageWithReply
	}

	return result, nil
}
