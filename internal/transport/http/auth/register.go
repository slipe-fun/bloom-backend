package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.KeysRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	pubKeys := req.EncryptedIdentityKeys.IdentityPublicKeys
	if pubKeys.MlKemPublicKey == "" || pubKeys.EcdhPublicKey == "" || pubKeys.EdPublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "missing required public keys",
		})
	}

	idSecret := req.EncryptedIdentityKeys.EncryptedSecretKeys
	if idSecret.Ciphertext == "" || idSecret.Nonce == "" || idSecret.Salt == "" || idSecret.Signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "missing required encrypted identity key fields",
		})
	}

	masterKey := req.EncryptedMasterKey
	if masterKey.Ciphertext == "" || masterKey.Nonce == "" || masterKey.Salt == "" || masterKey.Signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "missing required encrypted master key fields",
		})
	}

	token, user, session, err := h.authApp.Register(&req)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"token":   token,
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
