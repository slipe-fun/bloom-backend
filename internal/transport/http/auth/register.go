package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/config"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if !config.EmailRegex.MatchString(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_email",
			"message": "invalid email",
		})
	}

	err := h.authApp.Register(req.Email)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
