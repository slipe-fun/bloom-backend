package chat

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var req struct {
		Recipient int `json:"recipient"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	if req.Recipient == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_recipient",
			"message": "no recipient",
		})
	}

	user, err := h.userApp.GetUserByID(req.Recipient)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	chat, err := h.chatApp.GetChatWithUsers(session.UserID, req.Recipient)
	if chat != nil || err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "already_exists",
			"message": "chat with users already exists",
		})
	}

	chat, err = h.chatApp.CreateChat(session.UserID, user.ID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	outMsg := struct {
		Type   string `json:"type"`
		UserID int    `json:"user_id"`
		*domain.Chat
	}{
		Type:   "chat.new",
		UserID: session.UserID,
		Chat:   chat,
	}

	b, err := json.Marshal(outMsg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "internal error",
		})
	}

	h.wsHub.SendToUser(user.ID, b)

	return c.JSON(chat)
}
