package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/config"
)

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request"})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_request"})
	}

	if !config.EmailRegex.MatchString(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_email"})
	}

	err := h.authApp.Register(req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cant_create_user"})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
