package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	user, err := h.userApp.GetUserByID(id)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(fiber.Map{
		"id":           user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"description":  user.Description,
		"date":         user.Date,
	})
}
