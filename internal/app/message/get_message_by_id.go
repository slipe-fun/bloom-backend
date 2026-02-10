package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (m *MessageApp) GetMessageById(tokenStr string, id int) (*domain.MessageWithReply, error) {
	_, err := m.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	message, err := m.messages.GetById(id)
	if err != nil {
		logger.LogError(err.Error(), "message-app")
		return nil, domain.NotFound("message not found")
	}

	_, chaterr := m.chats.GetChatById(tokenStr, message.ChatID)
	if chaterr != nil {
		return nil, chaterr
	}

	result := &domain.MessageWithReply{
		Message: *message,
	}

	if message.ReplyTo != nil {
		replyToMessage, err := m.messages.GetById(*message.ReplyTo)
		if err == nil && replyToMessage != nil && replyToMessage.ChatID == message.ChatID {
			result.ReplyToMessage = replyToMessage
		}
	}

	return result, nil
}
