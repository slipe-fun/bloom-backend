package chat

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *ChatHandler) InviteUsersToGroup(c *fiber.Ctx) error {
	userVal := c.Locals("session_user")
	sessionUser, ok := userVal.(*domain.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	chatID, err := c.ParamsInt("c_id", 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request params",
		})
	}

	var req struct {
		Members []Member `json:"members"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request body",
		})
	}

	if chatID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "missing chat id",
		})
	}

	if len(req.Members) == 0 || len(req.Members) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid member count",
		})
	}

	chat, err := h.chatApp.GetChatByID(c.Context(), sessionUser.ID, chatID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": "chat not found",
		})
	}

	if chat.Type == "private" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "not a group",
		})
	}

	member, err := h.chatApp.GetGroupMember(c.Context(), sessionUser.ID, chatID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "not_found",
			"message": "chat not found",
		})
	}

	if member.Role != "creator" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "no permissions",
		})
	}

	groupMembers, buildErr := h.BuildGroupMemberList(MembersRequest{
		Members: req.Members,
	}, sessionUser)
	if buildErr != nil {
		return c.Status(buildErr.Status).JSON(buildErr.Response)
	}

	prevMembers, err := h.chatApp.GetGroupMembers(c.Context(), chat.ID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	users := make(map[int]domain.GroupMember, len(groupMembers))
	for i := range prevMembers {
		member := prevMembers[i]
		users[member.MemberID] = member
	}

	var resultMembers []domain.GroupMember
	for i := range groupMembers {
		member := groupMembers[i]

		_, exists := users[member.MemberID]
		if exists {
			continue
		}

		resultMembers = append(resultMembers, member)
	}

	if len(resultMembers) == 0 {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "all specified users are already members",
		})
	}

	err = h.chatApp.InviteUsersToGroup(c.Context(), chat.ID, sessionUser.ID, resultMembers)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	for i := range resultMembers {
		member := resultMembers[i]

		outMsg := struct {
			*domain.Chat
			Type                     string      `json:"type"`
			ChatType                 string      `json:"chat_type"`
			UserID                   string      `json:"user_id"`
			Role                     string      `json:"role"`
			InvitedBy                domain.User `json:"invited_by"`
			domain.Handshake         `json:"handshake"`
			domain.EncryptedGroupKey `json:"encrypted_group_key"`
		}{
			Chat:              chat,
			Type:              "chat.new",
			ChatType:          "group",
			UserID:            sessionUser.PublicID,
			Role:              "member",
			InvitedBy:         *sessionUser,
			Handshake:         member.Handshake,
			EncryptedGroupKey: member.EncryptedGroupKey,
		}

		b, err := json.Marshal(outMsg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "internal_error",
				"message": "internal error",
			})
		}

		h.wsHub.SendToUser(member.MemberID, b)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}
