package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) ExchangeCode(c *fiber.Ctx) error {
	if c.Query("state") != "random-state" {
		return c.Status(400).SendString("invalid state")
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(400).SendString("no code")
	}

	token, session, user, err := h.authApp.ExchangeCode(code)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(fiber.Map{
		"token":      token,
		"session_id": session.ID,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"date":     user.Date,
		},
	})
}
