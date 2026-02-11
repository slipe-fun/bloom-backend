package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *ChatHandler) GetChatMessagesAfter(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_token",
			"message": "invalid token",
		})
	}

	chatID, err := c.ParamsInt("c_id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	afterID, err := c.ParamsInt("m_id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	count := c.QueryInt("count")
	if count < 1 || count > 20 {
		count = 20
	}

	messages, err := h.messageApp.GetChatMessagesAfter(token, chatID, afterID, count)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	if len(messages) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no_messages",
		})
	}

	return c.JSON(messages)
}
