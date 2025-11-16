package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *SessionHandler) DeleteSession(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid_token",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_params",
		})
	}

	deleteSessionErr := h.sessionApp.DeleteSession(id, token)
	if deleteSessionErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed_to_delete_session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
