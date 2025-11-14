package auth

import (
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) RequestCode(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request"})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_request"})
	}

	if len(req.Email) < 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_username"})
	}

	err := h.authApp.RequestCode(req.Email)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cant_send_email"})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
