package keys

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *KeysHandler) GetUserChatsKeys(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid_token",
		})
	}

	keys, err := h.keysApp.GetUserChatsKeys(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cant_get_keys"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":    keys.UserID,
		"ciphertext": keys.Ciphertext,
		"nonce":      keys.Nonce,
		"salt":       keys.Salt,
	})
}
