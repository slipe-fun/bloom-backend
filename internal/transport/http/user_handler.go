package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/app"
)

type UserHandler struct {
	userApp *app.UserApp
}

func NewUserHandler(userApp *app.UserApp) *UserHandler {
	return &UserHandler{userApp: userApp}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	token, err := ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.userApp.GetUserByToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"date":     user.Date,
	})
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.userApp.GetUserById(id)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"date":     user.Date,
	})
}
