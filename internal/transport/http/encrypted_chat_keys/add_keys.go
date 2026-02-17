package encryptedchatkeys

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *EncryptedChatKeysHandler) AddKeys(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	chat_id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	var req struct {
		Keys []*domain.RawEncryptedChatKeys
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if len(req.Keys) == 0 || len(req.Keys) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	var keys []*domain.EncryptedChatKeys

	for _, key := range req.Keys {
		if key == nil ||
			key.SessionID <= 0 ||
			len(key.EncryptedKey) == 0 ||
			len(key.EncapsulatedKey) == 0 ||
			len(key.Nonce) == 0 ||
			len(key.Salt) == 0 {

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_request",
				"message": "invalid key payload",
			})
		}

		keys = append(keys, &domain.EncryptedChatKeys{
			ChatID:          chat_id,
			SessionID:       key.SessionID,
			EncryptedKey:    key.EncryptedKey,
			EncapsulatedKey: key.EncapsulatedKey,
			Nonce:           key.Nonce,
			Salt:            key.Salt,
		})
	}

	createdKeys, err := h.keys.AddKeys(session.UserID, chat_id, keys)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "something went wrong",
		})
	}

	return c.JSON(createdKeys)
}
