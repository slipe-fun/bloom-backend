package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *UserHandler) IsUserWithEmailExists(c *fiber.Ctx) error {
	email := c.Query("email", "")
	if email == "" {
		return c.JSON(fiber.Map{
			"exists":  false,
			"error":   "not_found",
			"message": "user not found",
		})
	}

	_, err := h.userApp.IsUserWithEmailExists(email)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"exists":  false,
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(fiber.Map{
		"exists": true,
	})
}
