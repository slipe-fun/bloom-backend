package auth

import (
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) VerifyCode(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_request"})
	}

	if req.Email == "" || req.Code == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_request"})
	}

	if len(req.Email) < 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_username"})
	}

	token, user, err := h.authApp.VerifyCode(req.Email, req.Code)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cant_login_user"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"date":     user.Date,
		},
	})
}
