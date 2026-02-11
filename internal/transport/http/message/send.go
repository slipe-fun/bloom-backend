package message

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *MessageHandler) Send(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var req struct {
		Type string `json:"type"`

		Ciphertext string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
		ChatID     int    `json:"chat_id"`
		ReplyTo    int    `json:"reply_to"`

		EncapsulatedKey       string `json:"encapsulated_key"`
		Signature             string `json:"signature"`
		SignedPayload         string `json:"signed_payload"`
		CEKWrap               string `json:"cek_wrap"`
		CEKWrapIV             string `json:"cek_wrap_iv"`
		CEKWrapSalt           string `json:"cek_wrap_salt"`
		EncapsulatedKeySender string `json:"encapsulated_key_sender"`
		CEKWrapSender         string `json:"cek_wrap_sender"`
		CEKWrapSenderIV       string `json:"cek_wrap_sender_iv"`
		CEKWrapSenderSalt     string `json:"cek_wrap_sender_salt"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.Type == "" || req.Ciphertext == "" || req.Nonce == "" || req.ChatID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	encryptionType := domain.EncryptionType(req.Type)
	if !encryptionType.IsValid() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if encryptionType == domain.ClientEncryption {
		if req.EncapsulatedKey == "" || req.Signature == "" || req.SignedPayload == "" ||
			req.CEKWrap == "" || req.CEKWrapIV == "" || req.CEKWrapSalt == "" ||
			req.EncapsulatedKeySender == "" || req.CEKWrapSender == "" || req.CEKWrapSenderIV == "" || req.CEKWrapSenderSalt == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_request",
				"message": "invalid request",
			})
		}
	}

	message, chat, err := h.messageApp.Send(session.UserID, encryptionType, &domain.SocketMessage{
		Ciphertext:            req.Ciphertext,
		Nonce:                 req.Nonce,
		ChatID:                req.ChatID,
		ReplyTo:               req.ReplyTo,
		EncapsulatedKey:       req.EncapsulatedKey,
		Signature:             req.Signature,
		SignedPayload:         req.SignedPayload,
		CEKWrap:               req.CEKWrap,
		CEKWrapIV:             req.CEKWrapIV,
		CEKWrapSalt:           req.CEKWrapSalt,
		EncapsulatedKeySender: req.EncapsulatedKeySender,
		CEKWrapSender:         req.CEKWrapSender,
		CEKWrapSenderIV:       req.CEKWrapSenderIV,
		CEKWrapSenderSalt:     req.CEKWrapSenderSalt,
	})
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	outMsg := struct {
		Type           string          `json:"type"`
		EncryptionType string          `json:"encryption_type"`
		ID             int             `json:"id"`
		UserID         int             `json:"user_id"`
		ReplyTo        *domain.Message `json:"reply_to,omitempty"`
		*domain.MessageWithReply
	}{
		Type:             "message.new",
		EncryptionType:   req.Type,
		ID:               message.ID,
		UserID:           session.UserID,
		ReplyTo:          message.ReplyToMessage,
		MessageWithReply: message,
	}

	b, err := json.Marshal(outMsg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "internal error",
		})
	}

	h.wsHub.SendToUser(h.chatApp.GetOtherMember(chat, session.UserID).ID, b)

	return c.JSON(message)
}
