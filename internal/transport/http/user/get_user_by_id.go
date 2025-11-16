package user

import "github.com/gofiber/fiber/v2"

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_params",
		})
	}

	user, err := h.userApp.GetUserById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "user_not_found",
		})
	}

	return c.JSON(fiber.Map{
		"id":           user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"date":         user.Date,
	})
}
