package encryptedchatkeys

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type addKeyBatchRequest struct {
	ChatID          int    `json:"chat_id"`
	RecipientID     int    `json:"recipient"`
	SessionID       int    `json:"session_id"`
	EncryptedKey    string `json:"encrypted_key"`
	EncapsulatedKey string `json:"encapsulated_key"`
	CekWrap         string `json:"cek_wrap"`
	CekWrapIV       string `json:"cek_wrap_iv"`
	Salt            string `json:"salt"`
	Nonce           string `json:"nonce"`
}

func (h *EncryptedChatKeysHandler) AddKeys(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var req []addKeyBatchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid request payload",
		})
	}

	if len(req) == 0 || len(req) > 200 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_request",
			"message": "invalid batch size",
		})
	}

	var batch []domain.AddKeyBatchItem

	sessionToRecipient := make(map[int]int, len(req))

	for _, key := range req {
		if key.ChatID <= 0 ||
			key.RecipientID <= 0 ||
			key.SessionID <= 0 ||
			len(key.EncryptedKey) == 0 ||
			len(key.EncapsulatedKey) == 0 ||
			len(key.Nonce) == 0 ||
			len(key.Salt) == 0 {

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_request",
				"message": "invalid key payload",
			})
		}

		sessionToRecipient[key.SessionID] = key.RecipientID

		domainKey := &domain.EncryptedChatKeys{
			ChatID:          key.ChatID,
			SessionID:       key.SessionID,
			FromSessionID:   session.ID,
			EncryptedKey:    key.EncryptedKey,
			EncapsulatedKey: key.EncapsulatedKey,
			CekWrap:         key.CekWrap,
			CekWrapIV:       key.CekWrapIV,
			Nonce:           key.Nonce,
			Salt:            key.Salt,
		}

		batch = append(batch, domain.AddKeyBatchItem{
			ChatID:      key.ChatID,
			RecipientID: key.RecipientID,
			Key:         domainKey,
		})
	}

	createdKeys, err := h.keys.AddKeys(session.UserID, batch)
	if err != nil {
		if appErr, ok := err.(*domain.AppError); ok {
			return c.Status(appErr.Status).JSON(fiber.Map{
				"error":   appErr.Code,
				"message": appErr.Msg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal_error",
			"message": "something went wrong",
		})
	}

	type wsMsg struct {
		Type          string                      `json:"type"`
		ChatID        int                         `json:"chat_id"`
		FromSessionID int                         `json:"from_session_id"`
		Keys          []*domain.EncryptedChatKeys `json:"keys"`
	}

	notifications := make(map[int]map[int][]*domain.EncryptedChatKeys)

	for _, k := range createdKeys {
		recID := sessionToRecipient[k.SessionID]
		if notifications[recID] == nil {
			notifications[recID] = make(map[int][]*domain.EncryptedChatKeys)
		}
		notifications[recID][k.ChatID] = append(notifications[recID][k.ChatID], k)
	}

	for recID, chats := range notifications {
		for chatID, keys := range chats {
			outMsg := wsMsg{
				Type:          "keys.new",
				ChatID:        chatID,
				FromSessionID: session.ID,
				Keys:          keys,
			}
			if b, err := json.Marshal(outMsg); err == nil {
				h.wsHub.SendToUser(recID, b)
			}
		}
	}

	return c.JSON(createdKeys)
}
