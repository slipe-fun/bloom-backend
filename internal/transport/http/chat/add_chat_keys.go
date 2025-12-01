package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *ChatHandler) AddChatKeys(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid_token",
		})
	}

	chatId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_params"})
	}

	if chatId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no_chat"})
	}

	var req struct {
		KyberPublicKey string `json:"kyberPublicKey"`
		EcdhPublicKey  string `json:"ecdhPublicKey"`
		EdPublicKey    string `json:"edPublicKey"`
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

	keysCheck := service.CheckKeysLength(req.KyberPublicKey, req.EcdhPublicKey, req.EdPublicKey)
	if keysCheck != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   keysCheck.Error(),
			"message": "invalid key length",
		})
	}

	chat, err := h.chatApp.GetChatById(token, chatId)
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

	return c.JSON(fiber.Map{
		"success": true,
	})
}
