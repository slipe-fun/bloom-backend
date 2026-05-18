package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
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

	return c.JSON(fiber.Map{
		"id":                user.ID,
		"session_id":        session.ID,
		"username":          user.Username,
		"display_name":      user.DisplayName,
		"description":       user.Description,
		"ml_kem_public_key": user.KyberPublicKey,
		"ecdh_public_key":   user.EcdhPublicKey,
		"ed_public_key":     user.EdPublicKey,
		"date":              user.Date,
	})
}
