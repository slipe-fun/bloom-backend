package MessageApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (c *MessageApp) GetMessageById(tokenStr string, id int) (*domain.MessageWithReply, error) {
	_, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	message, err := c.messages.GetById(id)
	if err != nil {
		return nil, domain.NotFound("message not found")
	}

	_, chaterr := c.chats.GetChatById(tokenStr, message.ChatID)
	if chaterr != nil {
		return nil, domain.NotFound("message not found")
	}

	result := &domain.MessageWithReply{
		Message: *message,
	}

	if message.ReplyTo != nil {
		replyToMessage, err := c.messages.GetById(*message.ReplyTo)
		if err == nil && replyToMessage != nil && replyToMessage.ChatID == message.ChatID {
			result.ReplyToMessage = replyToMessage
		}
	}

	return result, nil
}
