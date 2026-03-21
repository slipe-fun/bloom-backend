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
		Ciphertext string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
		ChatID     int    `json:"chat_id"`
		ReplyTo    int    `json:"reply_to"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.Ciphertext == "" || req.Nonce == "" || req.ChatID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	message, chat, err := h.messageApp.Send(session.UserID, &domain.SocketMessage{
		Ciphertext: req.Ciphertext,
		Nonce:      req.Nonce,
		ChatID:     req.ChatID,
		ReplyTo:    req.ReplyTo,
	})
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	outMsg := struct {
		Type    string          `json:"type"`
		ID      int             `json:"id"`
		UserID  int             `json:"user_id"`
		ReplyTo *domain.Message `json:"reply_to,omitempty"`
		*domain.MessageWithReply
	}{
		Type:             "message.new",
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
