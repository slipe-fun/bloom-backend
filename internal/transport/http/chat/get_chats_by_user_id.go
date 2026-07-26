package chat

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *ChatHandler) GetChatsByUserID(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	chats, err := h.chatApp.GetChatsByUserID(session.UserID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	groupsMember, err := h.chatApp.GetMemberGroups(c.Context(), session.UserID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	groupsMemberMap := make(map[int]domain.GroupMember, len(groupsMember))

	var usersIDs []int

	for i := range groupsMember {
		groupMember := groupsMember[i]
		groupsMemberMap[groupMember.ChatID] = groupMember
		usersIDs = append(usersIDs, groupMember.InvitedByID)
	}

	users, err := h.userApp.GetUsersByIDs(usersIDs)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "failed",
			"message": "failed to get users",
		})
	}

	usersMap := make(map[int]*domain.User, len(users))

	for i := range users {
		user := users[i]
		usersMap[user.ID] = &user
	}

	responseChats := make([]fiber.Map, 0, len(chats))

	for i := range chats {
		chat := chats[i]

		switch chat.Type {
		case "private":
			responseChats = append(responseChats, fiber.Map{
				"id":                chat.ID,
				"members":           chat.Members,
				"handshake":         chat.Handshake,
				"last_message":      chat.LastMessage,
				"last_read_message": chat.LastReadMessage,
			})
		case "group":
			member, exists := groupsMemberMap[chat.ID]
			if !exists {
				continue
			}

			var invitedBy *domain.User
			var handshake *domain.Handshake
			var encryptedGroupKey *domain.EncryptedGroupKey

			encryptedGroupKey = &member.EncryptedGroupKey
			handshake = &member.Handshake

			switch member.Role {
			case "creator":
				handshake = nil
				encryptedGroupKey = nil
			case "member":
				invitedBy = usersMap[member.InvitedByID]
			}

			responseChats = append(responseChats, fiber.Map{
				"id":                  chat.ID,
				"type":                chat.Type,
				"title":               chat.Title,
				"role":                member.Role,
				"invited_by":          invitedBy,
				"handshake":           handshake,
				"encrypted_group_key": encryptedGroupKey,
				"last_message":        chat.LastMessage,
				"last_read_message":   chat.LastReadMessage,
			})
		}

	}

	return c.JSON(responseChats)
}
