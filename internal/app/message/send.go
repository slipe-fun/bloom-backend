package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
	"github.com/slipe-fun/skid-backend/internal/pointer"
)

func (m *MessageApp) Send(user_id int, message *domain.SocketMessage) (*domain.MessageWithReply, *domain.Chat, error) {
	chat, err := m.chats.GetChatByID(user_id, message.ChatID)
	if err != nil {
		return nil, nil, domain.NotFound("chat not found")
	}

	var member *domain.Member
	for i, m := range chat.Members {
		if m.ID == user_id {
			member = &chat.Members[i]
			break
		}
	}
	if member == nil {
		return nil, nil, domain.NotFound("chat not found")
	}

	var replyTo *domain.Message
	if message.ReplyTo != 0 {
		reply_to_message, err := m.messages.GetByID(message.ReplyTo)
		if err != nil || reply_to_message == nil || reply_to_message.ChatID != chat.ID {
			return nil, nil, domain.NotFound("reply to message not found")
		}
		replyTo = reply_to_message
	}

	newMessage, err := m.messages.Create(&domain.Message{
		Ciphertext: message.Ciphertext,
		Nonce:      message.Nonce,
		ChatID:     message.ChatID,
		ReplyTo:    pointer.Intptr(message.ReplyTo),
	})
	if err != nil {
		logger.LogError(err.Error(), "message-app")
		return nil, nil, domain.Failed("failed to create message")
	}

	createdMessage := &domain.MessageWithReply{
		Message:        *newMessage,
		ReplyToMessage: replyTo,
	}

	return createdMessage, chat, nil
}
