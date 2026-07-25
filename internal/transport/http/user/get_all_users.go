package user

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ids := c.Query("ids")

	if ids == "" {
		limit := c.QueryInt("limit", 20)
		offset := c.QueryInt("offset", 0)

		if limit < 1 {
			limit = 20
		} else if limit > 20 {
			limit = 20
		}

		if offset < 0 {
			offset = 0
		}

		users, err := h.userApp.GetAllUsers(limit, offset)
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		return c.JSON(users)
	}

	usersIDs := strings.Split(ids, ",")

	if len(usersIDs) > 100 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": "too many users",
		})
	}

	users, err := h.userApp.GetUsersByPublicIDs(usersIDs)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(users)
}
