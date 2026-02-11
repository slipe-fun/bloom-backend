package chat

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto/validations"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *ChatHandler) AddChatKeys(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid_token",
		})
	}

	session, err := h.wsHub.SessionApp.GetSession(token)
	if err != nil {
		return err
	}

	chatID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_params"})
	}

	if chatID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no_chat"})
	}

	var req struct {
		KyberPublicKey string `json:"kyber_public_key"`
		EcdhPublicKey  string `json:"ecdh_public_key"`
		EdPublicKey    string `json:"ed_public_key"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.KyberPublicKey == "" || req.EcdhPublicKey == "" || req.EdPublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "not_all_keys",
			"message": "not all keys are provided",
		})
	}

	keysCheck := validations.CheckKeysLength(req.KyberPublicKey, req.EcdhPublicKey, req.EdPublicKey)
	if keysCheck != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   keysCheck.Error(),
			"message": "invalid key length",
		})
	}

	chat, err := h.chatApp.GetChatByID(token, chatID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	updateChatErr := h.chatApp.AddKeys(token, chat, req.KyberPublicKey, req.EcdhPublicKey, req.EdPublicKey)
	if appErr, ok := updateChatErr.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	outMsg := struct {
		Type   string `json:"type"`
		UserID int    `json:"user_id"`
		ChatID int    `json:"chat_id"`
		domain.SocketKeys
	}{
		Type:   "chat.keys_updated",
		UserID: session.UserID,
		ChatID: chat.ID,
		SocketKeys: domain.SocketKeys{
			ChatID:         chat.ID,
			KyberPublicKey: req.KyberPublicKey,
			EcdhPublicKey:  req.EcdhPublicKey,
			EdPublicKey:    req.EdPublicKey,
		},
	}

	b, err := json.Marshal(outMsg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "internal error",
		})
	}

	h.wsHub.SendToUser(h.chatApp.GetOtherMember(chat, session.UserID).ID, b)

	return c.JSON(fiber.Map{
		"success": true,
	})
}
