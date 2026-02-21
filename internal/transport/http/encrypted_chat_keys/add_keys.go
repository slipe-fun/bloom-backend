package encryptedchatkeys

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *EncryptedChatKeysHandler) AddKeys(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	chatID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	var req []*domain.RawEncryptedChatKeys
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if len(req) == 0 || len(req) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	var keys []*domain.EncryptedChatKeys

	for _, key := range req {
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
			ChatID:          chatID,
			SessionID:       key.SessionID,
			EncryptedKey:    key.EncryptedKey,
			EncapsulatedKey: key.EncapsulatedKey,
			CekWrap:         key.CekWrap,
			CekWrapIV:       key.CekWrapIV,
			Nonce:           key.Nonce,
			Salt:            key.Salt,
		})
	}

	createdKeys, recipientID, err := h.keys.AddKeys(session.UserID, chatID, keys)
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

	outMsg := struct {
		Type   string                      `json:"type"`
		ChatID int                         `json:"chat_id"`
		Keys   []*domain.EncryptedChatKeys `json:"keys"`
	}{
		Type:   "keys.new",
		ChatID: chatID,
		Keys:   createdKeys,
	}

	b, err := json.Marshal(outMsg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "internal error",
		})
	}

	h.wsHub.SendToUser(recipientID, b)

	return c.JSON(createdKeys)
}
