package message

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *MessageHandler) GetMessageByID(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_token",
			"message": "invalid token",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	message, err := h.messageApp.GetMessageByID(token, id)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(message)
}
