package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) RegisterBegin(c *fiber.Ctx) error {
	token, username, options, err := h.authApp.BeginRegistration()
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token":    token,
		"username": username,
		"options":  options,
	})
}
