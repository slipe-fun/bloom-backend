package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (m *MessageApp) GetChatMessages(token string, chatID int) ([]*domain.MessageWithReply, error) {
	session, err := m.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	chat, err := m.chats.GetChatByID(session.UserID, chatID)
	if err != nil {
		return nil, err
	}

	messages, err := m.messages.GetChatMessages(chat.ID)
	if err != nil {
		logger.LogError(err.Error(), "message-app")
		return nil, domain.NotFound("messages not found")
	}

	result := make([]*domain.MessageWithReply, len(messages))
	for i, msg := range messages {
		messageWithReply := &domain.MessageWithReply{
			Message: *msg,
		}

		if msg.ReplyTo != nil {
			replyToMessage, err := m.messages.GetByID(*msg.ReplyTo)
			if err == nil && replyToMessage != nil && replyToMessage.ChatID == msg.ChatID {
				messageWithReply.ReplyToMessage = replyToMessage
			}
		}

		result[i] = messageWithReply
	}

	return result, nil
}
