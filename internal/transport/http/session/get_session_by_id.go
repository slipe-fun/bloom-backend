package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *SessionHandler) GetSessionByToken(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid_token",
		})
	}

	session, err := h.sessionApp.GetSession(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed_to_get_session",
		})
	}

	return c.JSON(session)
}
