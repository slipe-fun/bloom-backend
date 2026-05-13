package keys

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *KeysHandler) GetUserKeys(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	keys_type := c.Params("type")
	if !(keys_type == "master" || keys_type == "bundle") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "wrong type",
		})
	}

	keys, err := h.keysApp.GetUserKeys(session.UserID, keys_type)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":    keys.UserID,
		"ciphertext": keys.Ciphertext,
		"nonce":      keys.Nonce,
		"salt":       keys.Salt,
	})
}
