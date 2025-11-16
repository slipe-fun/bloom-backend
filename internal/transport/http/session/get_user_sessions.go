package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *SessionHandler) GetUserSessions(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid_token",
		})
	}

	sessions, err := h.sessionApp.GetUserSessions(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed_to_get_sessions",
		})
	}

	return c.JSON(sessions)
}
