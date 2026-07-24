package exchange

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *ExchangeHandler) StartSession(c *fiber.Ctx) error {
	roomID, err := h.exchangeApp.StartSession()
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(fiber.Map{
		"room_id": roomID,
	})
}
