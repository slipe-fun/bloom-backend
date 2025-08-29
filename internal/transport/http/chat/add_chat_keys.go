package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/service"
	"github.com/slipe-fun/skid-backend/internal/transport/http"
)

func (h *ChatHandler) AddChatKeys(c *fiber.Ctx) error {
	token, err := http.ExtractBearerToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	chatId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if chatId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "chat id required"})
	}

	var req struct {
		KyberPublicKey string `json:"kyberPublicKey"`
		EcdhPublicKey  string `json:"ecdhPublicKey"`
		EdPublicKey    string `json:"edPublicKey"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.KyberPublicKey == "" || req.EcdhPublicKey == "" || req.EdPublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "u need to specify all keys"})
	}

	keysCheck := service.CheckKeysLength(req.KyberPublicKey, req.EcdhPublicKey, req.EdPublicKey)
	if keysCheck != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": keysCheck.Error()})
	}

	chat, err := h.chatApp.GetChatById(token, chatId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updateChatErr := h.chatApp.AddKeys(token, chat, req.KyberPublicKey, req.EcdhPublicKey, req.EdPublicKey)
	if updateChatErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": updateChatErr.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
