package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *UserHandler) EditUser(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user, err := h.userApp.GetUserByID(session.UserID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	var req struct {
		Username    *string `json:"username"`
		DisplayName *string `json:"display_name" form:"display_name"`
		Description *string `json:"description" form:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.Username == nil && req.DisplayName == nil && req.Description == nil {
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

	if req.Description != nil {
		newUser.Description = req.Description
	} else {
		newUser.Description = user.Description
	}

	edited, err := h.userApp.EditUser(session.UserID, newUser)
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
			"description":  edited.Description,
			"date":         edited.Date,
		},
	})
}
