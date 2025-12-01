package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *UserHandler) EditUser(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid_token",
			"message": "invalid token",
		})
	}

	user, err := h.userApp.GetUserByToken(token)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	var req struct {
		Username    *string `json:"username"`
		DisplayName *string `json:"display_name" form:"display_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	newUser := &domain.User{
		ID:    user.ID,
		Email: user.Email,
		Date:  user.Date,
	}

	if req.Username != nil {
		newUser.Username = *req.Username
	} else {
		newUser.Username = user.Username
	}

	if req.DisplayName != nil {
		newUser.DisplayName = req.DisplayName
	} else {
		newUser.DisplayName = user.DisplayName
	}

	edited, err := h.userApp.EditUser(token, newUser)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":           edited.ID,
			"username":     edited.Username,
			"display_name": edited.DisplayName,
			"date":         edited.Date,
		},
	})
}
