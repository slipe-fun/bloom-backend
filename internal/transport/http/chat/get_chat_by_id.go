package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *ChatHandler) GetChatByID(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	chat, err := h.chatApp.GetChatByID(c.Context(), session.UserID, id)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	if chat.Type == "group" {
		member, err := h.chatApp.GetGroupMember(c.Context(), chat.ID, session.UserID)
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		invitedBy, err := h.userApp.GetUserByID(member.InvitedByID)
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		var handshake *domain.Handshake
		var encryptedGroupKey *domain.EncryptedGroupKey

		encryptedGroupKey = &member.EncryptedGroupKey
		handshake = &member.Handshake

		if member.Role == "creator" {
			handshake = nil
			encryptedGroupKey = nil
			invitedBy = nil
		}

		return c.JSON(fiber.Map{
			"id":                  chat.ID,
			"type":                chat.Type,
			"title":               chat.Title,
			"role":                member.Role,
			"invited_by":          invitedBy,
			"handshake":           handshake,
			"encrypted_group_key": encryptedGroupKey,
		})
	}

	return c.JSON(fiber.Map{
		"id":      chat.ID,
		"members": chat.Members,
	})
}
