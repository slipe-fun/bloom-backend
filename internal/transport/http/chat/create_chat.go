package chat

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto/validations"
)

const (
	MLKEM768CiphertextSize = 1088
	AESGCMNonceSize        = 12
	SyncKeyCiphertextSize  = 48
)

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var req struct {
		Recipient int `json:"recipient"`
		Handshake struct {
			ReceiverCipherText string `json:"receiver_cipher_text"`
			SenderCipherText   string `json:"sender_cipher_text"`
			EncryptedSyncKey   struct {
				CipherText string `json:"ciphertext"`
				Nonce      string `json:"nonce"`
			} `json:"encrypted_sync_key"`
		} `json:"handshake"`
	}

	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.Recipient == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_recipient",
			"message": "no recipient",
		})
	}

	if req.Handshake.ReceiverCipherText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_receiver_cipher_text",
			"message": "receiver cipher text is missing",
		})
	}

	if req.Handshake.SenderCipherText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_sender_cipher_text",
			"message": "sender cipher text is missing",
		})
	}

	if req.Handshake.EncryptedSyncKey.CipherText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_sync_key_ciphertext",
			"message": "encrypted sync key ciphertext is missing",
		})
	}

	if req.Handshake.EncryptedSyncKey.Nonce == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_sync_key_nonce",
			"message": "encrypted sync key nonce is missing",
		})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.ReceiverCipherText, MLKEM768CiphertextSize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_receiver_ciphertext_length"})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.SenderCipherText, MLKEM768CiphertextSize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_sender_ciphertext_length"})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.EncryptedSyncKey.Nonce, AESGCMNonceSize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_nonce_length"})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.EncryptedSyncKey.CipherText, SyncKeyCiphertextSize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_sync_key_length"})
	}

	if req.Recipient == session.UserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "self_chat_not_allowed",
		})
	}

	user, err := h.userApp.GetUserByID(req.Recipient)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	chat, err := h.chatApp.GetChatWithUsers(session.UserID, req.Recipient)
	if chat != nil || err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "already_exists",
			"message": "chat with users already exists",
		})
	}

	chat, err = h.chatApp.CreateChat(session.UserID, user.ID, domain.Handshake{
		ReceiverCipherText: req.Handshake.ReceiverCipherText,
		SenderCipherText:   req.Handshake.SenderCipherText,
		EncryptedSyncKey: struct {
			CipherText string `json:"ciphertext"`
			Nonce      string `json:"nonce"`
		}{
			CipherText: req.Handshake.EncryptedSyncKey.CipherText,
			Nonce:      req.Handshake.EncryptedSyncKey.Nonce,
		},
	})
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	outMsg := struct {
		Type   string `json:"type"`
		UserID int    `json:"user_id"`
		*domain.Chat
	}{
		Type:   "chat.new",
		UserID: session.UserID,
		Chat:   chat,
	}

	b, err := json.Marshal(outMsg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "internal error",
		})
	}

	h.wsHub.SendToUser(user.ID, b)

	return c.JSON(chat)
}
