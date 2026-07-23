package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) LoginBegin(c *fiber.Ctx) error {
	authLookupID := c.Params("auth_lookup_id")

	keys, challenge, userID, err := h.authApp.LoginBegin(authLookupID)
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
		"user_id":   userID,
		"keys":      keys,
		"challenge": challenge,
	})
}
