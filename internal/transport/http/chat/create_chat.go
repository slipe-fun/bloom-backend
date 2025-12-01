package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_token",
			"message": "invalid token",
		})
	}

	var req struct {
		Recipient int `json:"recipient"`
	}

	if err := c.BodyParser(&req); err != nil {
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

	user, err := h.userApp.GetUserById(req.Recipient)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	chat1, err := h.chatApp.GetChatWithUsers(token, req.Recipient)
	if chat1 != nil || err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "already_exists",
			"message": "chat with users already exists",
		})
	}

	chat, err := h.chatApp.CreateChat(token, user.ID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(fiber.Map{
		"id":             chat.ID,
		"members":        chat.Members,
		"encryption_key": chat.EncryptionKey,
	})
}
