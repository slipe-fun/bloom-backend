package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto/validations"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
	"github.com/slipe-fun/skid-backend/internal/pointer"
)

func (m *MessageApp) Send(user_id int, encryptionType domain.EncryptionType, message *domain.SocketMessage) (*domain.MessageWithReply, *domain.Chat, error) {
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

	var createdMessage *domain.MessageWithReply
	switch encryptionType {
	case domain.ServerEncryption:
		message, err := m.messages.Create(&domain.Message{
			Ciphertext: message.Ciphertext,
			Nonce:      message.Nonce,
			ChatID:     message.ChatID,
			ReplyTo:    pointer.Intptr(message.ReplyTo),
		})
		if err != nil {
			logger.LogError(err.Error(), "message-app")
			return nil, nil, domain.Failed("failed to create message")
		}

		createdMessage = &domain.MessageWithReply{
			Message:        *message,
			ReplyToMessage: replyTo,
		}
	case domain.ClientEncryption:
		if err := crypto.VerifySignature(
			member.EdPublicKey,
			message.SignedPayload,
			message.Signature,
		); err != nil {
			logger.LogError("invalid message signature", "message-app")
			return nil, nil, domain.Failed("invalid message signature")
		}

		if err := validations.ValidateCEKFields(
			message.CEKWrap,
			message.CEKWrapIV,
			message.CEKWrapSalt,
			message.EncapsulatedKeySender,
			message.CEKWrapSender,
			message.CEKWrapSenderIV,
			message.CEKWrapSenderSalt,
		); err != nil {
			return nil, nil, domain.Failed("invalid CEK fields")
		}

		message, err := m.messages.Create(&domain.Message{
			Ciphertext: message.Ciphertext,
			Nonce:      message.Nonce,
			ChatID:     message.ChatID,
			ReplyTo:    pointer.Intptr(message.ReplyTo),

			EncapsulatedKey:       pointer.Strptr(message.EncapsulatedKey),
			Signature:             pointer.Strptr(message.Signature),
			SignedPayload:         pointer.Strptr(message.SignedPayload),
			CEKWrap:               pointer.Strptr(message.CEKWrap),
			CEKWrapIV:             pointer.Strptr(message.CEKWrapIV),
			CEKWrapSalt:           pointer.Strptr(message.CEKWrapSalt),
			EncapsulatedKeySender: pointer.Strptr(message.EncapsulatedKeySender),
			CEKWrapSender:         pointer.Strptr(message.CEKWrapSender),
			CEKWrapSenderIV:       pointer.Strptr(message.CEKWrapSenderIV),
			CEKWrapSenderSalt:     pointer.Strptr(message.CEKWrapSenderSalt),
		})
		if err != nil {
			logger.LogError(err.Error(), "message-app")
			return nil, nil, domain.Failed("failed to create message")
		}

		createdMessage = &domain.MessageWithReply{
			Message:        *message,
			ReplyToMessage: replyTo,
		}
	default:
		return nil, nil, domain.Failed("unsupported encryption type")
	}

	return createdMessage, chat, nil
}
