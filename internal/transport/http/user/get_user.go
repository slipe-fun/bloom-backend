package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid_token",
		})
	}

	user, err := h.userApp.GetUserByToken(token)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user_not_found",
		})
	}

	return c.JSON(fiber.Map{
		"id":           user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"date":         user.Date,
	})
}
