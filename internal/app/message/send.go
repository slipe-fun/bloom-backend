package MessageApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/service/crypto"
	"github.com/slipe-fun/skid-backend/internal/service/logger"
)

func (m *MessageApp) Send(token, encryptionType string, message *domain.SocketMessage) (*domain.MessageWithReply, *domain.Chat, *domain.Session, error) {
	session, err := m.sessionApp.GetSession(token)
	if err != nil {
		return nil, nil, nil, err
	}

	chat, err := m.chats.GetChatById(token, message.ChatID)
	if err != nil {
		return nil, nil, nil, domain.NotFound("chat not found")
	}

	var member *domain.Member
	for i, m := range chat.Members {
		if m.ID == session.UserID {
			member = &chat.Members[i]
			break
		}
	}
	if member == nil {
		return nil, nil, nil, domain.NotFound("chat not found")
	}

	var replyTo *domain.Message
	if message.ReplyTo != 0 {
		reply_to_message, err := m.messages.GetById(message.ReplyTo)
		if err != nil || reply_to_message == nil || reply_to_message.ChatID != chat.ID {
			return nil, nil, nil, domain.NotFound("reply to message not found")
		}
		replyTo = reply_to_message
	}

	var createdMessage *domain.MessageWithReply
	switch encryptionType {
	case "server":
		message, err := m.messages.Create(&domain.Message{
			Ciphertext: message.Ciphertext,
			Nonce:      message.Nonce,
			ChatID:     message.ChatID,
			ReplyTo:    service.Intptr(message.ReplyTo),
		})
		if err != nil {
			logger.LogError(err.Error(), "message-app")
			return nil, nil, nil, domain.Failed("failed to create message")
		}

		createdMessage = &domain.MessageWithReply{
			Message:        *message,
			ReplyToMessage: replyTo,
		}
	case "client":
		if err := crypto.VerifySignature(
			member.EdPublicKey,
			message.SignedPayload,
			message.Signature,
		); err != nil {
			logger.LogError("invalid message signature", "message-app")
			return nil, nil, nil, domain.Failed("invalid message signature")
		}

		if err := crypto.ValidateCEKFields(
			message.CEKWrap,
			message.CEKWrapIV,
			message.CEKWrapSalt,
			message.EncapsulatedKeySender,
			message.CEKWrapSender,
			message.CEKWrapSenderIV,
			message.CEKWrapSenderSalt,
		); err != nil {
			return nil, nil, nil, domain.Failed("invalid CEK fields")
		}

		message, err := m.messages.Create(&domain.Message{
			Ciphertext: message.Ciphertext,
			Nonce:      message.Nonce,
			ChatID:     message.ChatID,
			ReplyTo:    service.Intptr(message.ReplyTo),

			EncapsulatedKey:       service.Strptr(message.EncapsulatedKey),
			Signature:             service.Strptr(message.Signature),
			SignedPayload:         service.Strptr(message.SignedPayload),
			CEKWrap:               service.Strptr(message.CEKWrap),
			CEKWrapIV:             service.Strptr(message.CEKWrapIV),
			CEKWrapSalt:           service.Strptr(message.CEKWrapSalt),
			EncapsulatedKeySender: service.Strptr(message.EncapsulatedKeySender),
			CEKWrapSender:         service.Strptr(message.CEKWrapSender),
			CEKWrapSenderIV:       service.Strptr(message.CEKWrapSenderIV),
			CEKWrapSenderSalt:     service.Strptr(message.CEKWrapSenderSalt),
		})
		if err != nil {
			logger.LogError(err.Error(), "message-app")
			return nil, nil, nil, domain.Failed("failed to create message")
		}

		createdMessage = &domain.MessageWithReply{
			Message:        *message,
			ReplyToMessage: replyTo,
		}
	default:
		return nil, nil, nil, domain.Failed("unsupported encryption type")
	}

	return createdMessage, chat, session, nil
}
