package chat

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto/validations"
)

const (
	MLKEM768CiphertextSize = 1088
	AESGCMNonceSize        = 12
	SyncKeyCiphertextSize  = 172
	X448PublicKeySize      = 56
)

func (h *ChatHandler) CreateChat(c *fiber.Ctx) error {
	userVal := c.Locals("session_user")
	sessionUser, ok := userVal.(*domain.User)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var req struct {
		Type string `json:"type"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request",
		})
	}

	switch req.Type {
	case "private":
		return h.createPrivateChat(c, sessionUser)
	case "group":
		return h.createGroupChat(c, sessionUser)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid chat type",
		})
	}
}

func (h *ChatHandler) createGroupChat(c *fiber.Ctx, sessionUser *domain.User) error {
	var req struct {
		Title   string `json:"title"`
		Members []struct {
			MemberID          string                   `json:"member_id"`
			Handshake         domain.Handshake         `json:"handshake"`
			EncryptedGroupKey domain.EncryptedGroupKey `json:"encrypted_group_key"`
		} `json:"members"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request body",
		})
	}

	if len(req.Title) < 1 || len(req.Title) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid title length",
		})
	}

	var members []domain.User
	users := make(map[string]domain.User, len(req.Members))
	var err error
	if len(req.Members) != 0 {
		var membersIDs []string
		for i := range req.Members {
			member := req.Members[i]
			membersIDs = append(membersIDs, member.MemberID)
		}

		members, err = h.userApp.GetUsersByPublicIDs(membersIDs)
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		if len(members) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_request",
				"message": "invalid members list",
			})
		}

		for i := range members {
			member := members[i]
			users[member.PublicID] = member
		}
	}

	var groupMembers []domain.GroupMember

	for i := range req.Members {
		memberObject := req.Members[i]

		member, exists := users[memberObject.MemberID]
		if !exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "user_not_found",
				"message": "member " + memberObject.MemberID + " not found",
			})
		}

		if member.ID == sessionUser.ID {
			continue
		}

		if memberObject.Handshake.ReceiverCipherText == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "no_receiver_cipher_text",
				"message": memberObject.MemberID + ": receiver cipher text is missing",
			})
		}

		if memberObject.Handshake.SenderEphemeralX448 == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "no_sender_cipher_text",
				"message": memberObject.MemberID + ": sender cipher text is missing",
			})
		}

		if memberObject.Handshake.EncryptedSyncKey.CipherText == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "no_sync_key_ciphertext",
				"message": memberObject.MemberID + ": encrypted sync key ciphertext is missing",
			})
		}

		if memberObject.Handshake.EncryptedSyncKey.Nonce == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "no_sync_key_nonce",
				"message": memberObject.MemberID + ": encrypted sync key nonce is missing",
			})
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.ReceiverCipherText, MLKEM768CiphertextSize); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "invalid_receiver_ciphertext_length",
				"message": memberObject.MemberID + ": invalid receiver ciphertext length",
			})
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.SenderEphemeralX448, X448PublicKeySize); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "invalid_sender_ciphertext_length",
				"message": memberObject.MemberID + ": invalid sender ciphertext length",
			})
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.EncryptedSyncKey.Nonce, AESGCMNonceSize); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "invalid_nonce_length",
				"message": memberObject.MemberID + ": invalid nonce length",
			})
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.EncryptedSyncKey.CipherText, SyncKeyCiphertextSize); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "invalid_sync_key_length",
				"message": memberObject.MemberID + ": invalid sync key length",
			})
		}

		newGroupMember := domain.GroupMember{
			MemberID:          member.ID,
			Role:              "member",
			Handshake:         memberObject.Handshake,
			EncryptedGroupKey: memberObject.EncryptedGroupKey,
		}

		groupMembers = append(groupMembers, newGroupMember)
	}

	creatorGroupMember := domain.GroupMember{
		MemberID: sessionUser.ID,
		Role:     "creator",
	}
	groupMembers = append(groupMembers, creatorGroupMember)

	chat, err := h.chatApp.CreateGroupChat(c.Context(), req.Title, groupMembers, sessionUser.ID)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	for i := range groupMembers {
		member := groupMembers[i]

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
		"id":    chat.ID,
		"type":  chat.Type,
		"title": chat.Title,
		"role":  "creator",
	})
}

func (h *ChatHandler) createPrivateChat(c *fiber.Ctx, sessionUser *domain.User) error {
	var req struct {
		Recipient string `json:"recipient"`
		Handshake struct {
			ReceiverCipherText  string `json:"receiver_cipher_text"`
			SenderEphemeralX448 string `json:"sender_ephemeral_x448"`
			EncryptedSyncKey    struct {
				CipherText string `json:"ciphertext"`
				Nonce      string `json:"nonce"`
			} `json:"encrypted_sync_key"`
		} `json:"handshake"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request body",
		})
	}

	if req.Recipient == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_recipient",
			"message": "no recipient",
		})
	}

	if req.Handshake.ReceiverCipherText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_receiver_cipher_text",
			"message": "receiver cipher text is missing",
		})
	}

	if req.Handshake.SenderEphemeralX448 == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_sender_cipher_text",
			"message": "sender cipher text is missing",
		})
	}

	if req.Handshake.EncryptedSyncKey.CipherText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_sync_key_ciphertext",
			"message": "encrypted sync key ciphertext is missing",
		})
	}

	if req.Handshake.EncryptedSyncKey.Nonce == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "no_sync_key_nonce",
			"message": "encrypted sync key nonce is missing",
		})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.ReceiverCipherText, MLKEM768CiphertextSize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_receiver_ciphertext_length"})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.SenderEphemeralX448, X448PublicKeySize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_sender_ciphertext_length"})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.EncryptedSyncKey.Nonce, AESGCMNonceSize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_nonce_length"})
	}

	if err := validations.ValidateCryptoLength(req.Handshake.EncryptedSyncKey.CipherText, SyncKeyCiphertextSize); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_sync_key_length"})
	}

	user, err := h.userApp.GetUserByPublicID(req.Recipient)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	if user.ID == sessionUser.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "self_chat_not_allowed",
		})
	}

	chat, err := h.chatApp.GetChatWithUsers(sessionUser.ID, user.ID)
	if chat != nil || err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "already_exists",
			"message": "chat with users already exists",
		})
	}

	chat, err = h.chatApp.CreatePrivateChat(sessionUser.ID, user.ID, domain.Handshake{
		ReceiverCipherText:  req.Handshake.ReceiverCipherText,
		SenderEphemeralX448: req.Handshake.SenderEphemeralX448,
		EncryptedSyncKey: domain.EncryptedSyncKey{
			CipherText: req.Handshake.EncryptedSyncKey.CipherText,
			Nonce:      req.Handshake.EncryptedSyncKey.Nonce,
		},
	})
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	outMsg := struct {
		*domain.Chat
		Type     string `json:"type"`
		ChatType string `json:"chat_type"`
		UserID   string `json:"user_id"`
	}{
		Chat:     chat,
		Type:     "chat.new",
		ChatType: "private",
		UserID:   sessionUser.PublicID,
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
