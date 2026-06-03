package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) LoginFinish(c *fiber.Ctx) error {
	var req struct {
		UserID    string `json:"user_id"`
		Signature string `json:"signature"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request body format",
		})
	}

	if req.UserID == "" || req.Signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "user_id and signature are required",
		})
	}

	token, user, session, err := h.authApp.LoginFinish(req.UserID, req.Signature)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "failed to finish login",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":                user.PublicID,
			"username":          user.Username,
			"display_name":      user.DisplayName,
			"description":       user.Description,
			"ml_kem_public_key": user.KyberPublicKey,
			"ecdh_public_key":   user.EcdhPublicKey,
			"ed_public_key":     user.EdPublicKey,
			"date":              user.Date,
		},
		"session": fiber.Map{
			"id":         session.ID,
			"token":      session.Token,
			"user_id":    user.PublicID,
			"revoked_at": session.RevokedAt,
			"created_at": session.CreatedAt,
		},
	})
}
