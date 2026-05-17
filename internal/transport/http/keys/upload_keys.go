package keys

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *KeysHandler) SaveKeys(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	keys_type := c.Params("type")
	if !(keys_type == "master") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "wrong type",
		})
	}

	var req struct {
		Ciphertext string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
		Salt       string `json:"salt"`
		Signature  string `json:"signature"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.Ciphertext == "" || req.Nonce == "" || req.Salt == "" || req.Signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	_, err := h.keysApp.CreateKeys(session.UserID, &domain.EncryptedKeys{
		Type:       keys_type,
		Ciphertext: req.Ciphertext,
		Nonce:      req.Nonce,
		Salt:       req.Salt,
		Signature:  req.Signature,
	})
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
