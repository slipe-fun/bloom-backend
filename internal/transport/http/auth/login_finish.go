package auth

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) LoginFinish(c *fiber.Ctx) error {
	var req struct {
		Token    string      `json:"token"`
		Response interface{} `json:"response"`
	}

	if err := c.BodyParser(&req); err != nil || req.Token == "" || req.Response == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid token or response payload",
		})
	}

	respBytes, err := json.Marshal(req.Response)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "failed to marshal response payload",
		})
	}

	sessionToken, session, user, err := h.authApp.FinishLogin(req.Token, respBytes)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": sessionToken,
		"session": fiber.Map{
			"id":         session.ID,
			"user_id":    session.UserID,
			"created_at": session.CreatedAt,
		},
		"user": fiber.Map{
			"id":           user.ID,
			"username":     user.Username,
			"display_name": user.DisplayName,
			"description":  user.Description,
			"date":         user.Date,
		},
	})
}
