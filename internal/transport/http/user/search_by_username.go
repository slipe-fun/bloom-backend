package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *UserHandler) SearchByUsername(c *fiber.Ctx) error {
	query := c.Query("q", "")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	if len(query) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": "users not found",
		})
	}

	users, err := h.userApp.SearchUsersByUsername(query, limit, offset)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(users)
}
