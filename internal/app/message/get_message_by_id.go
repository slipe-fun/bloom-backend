package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (m *MessageApp) GetMessageByID(token string, id int) (*domain.MessageWithReply, error) {
	_, err := m.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	message, err := m.messages.GetByID(id)
	if err != nil {
		logger.LogError(err.Error(), "message-app")
		return nil, domain.NotFound("message not found")
	}

	_, chaterr := m.chats.GetChatByID(token, message.ChatID)
	if chaterr != nil {
		return nil, chaterr
	}

	result := &domain.MessageWithReply{
		Message: *message,
	}

	if message.ReplyTo != nil {
		replyToMessage, err := m.messages.GetByID(*message.ReplyTo)
		if err == nil && replyToMessage != nil && replyToMessage.ChatID == message.ChatID {
			result.ReplyToMessage = replyToMessage
		}
	}

	return result, nil
}
