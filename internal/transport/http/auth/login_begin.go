package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) LoginBegin(c *fiber.Ctx) error {
	user_id := c.Params("user_id")

	keys, challenge, err := h.authApp.LoginBegin(user_id)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "failed to begin login",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"keys":      keys,
		"challenge": challenge,
	})
}
