package chat

import (
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) GetByUserID(userID int) ([]*domain.ChatWithLastMessage, error) {
	start := time.Now()

	rows, err := r.db.Query(`
			SELECT id, members, handshake, title, type
			FROM chats
			WHERE
				(type = 'private' AND EXISTS (
					SELECT 1
					FROM jsonb_array_elements(members) m
					WHERE (m->>'id')::int = $1 -- Приводим JSON-строку к числу и сравниваем инт с интом
				))
				OR
				(type = 'group' AND EXISTS (
					SELECT 1
					FROM group_members gm
					WHERE gm.chat_id = chats.id AND gm.member_id = $1 -- Сравниваем инт с интом
				))
		`, userID)

	duration := time.Since(start)

	metrics.ObserveDB("chats_get_with_user_id", duration, err)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*domain.ChatWithLastMessage

	for rows.Next() {
		var chat domain.ChatWithLastMessage
		var membersJSON []byte
		var handshakeJSON []byte
		var titlePtr *string

		if err := rows.Scan(&chat.ID, &membersJSON, &handshakeJSON, &titlePtr, &chat.Type); err != nil {
			return nil, err
		}

		if titlePtr != nil {
			chat.Title = *titlePtr
		}

		if chat.Type == "private" {
			var rawMembers []domain.Member
			if len(membersJSON) > 0 {
				if err := json.Unmarshal(membersJSON, &rawMembers); err != nil {
					return nil, err
				}
			}

			chat.Members = make([]domain.User, 0, len(rawMembers))

			for _, m := range rawMembers {
				user, err := r.userRepo.GetByID(m.ID)
				if err != nil {
					continue
				}
				chat.Members = append(chat.Members, *user)
			}

			if len(handshakeJSON) > 0 {
				var hs domain.Handshake
				if err := json.Unmarshal(handshakeJSON, &hs); err == nil {
					chat.Handshake = &hs
				}
			}
		}

		var msg domain.Message
		err = r.db.QueryRow(`
			SELECT id, ciphertext, nonce, chat_id, seen, reply_to
			FROM messages
			WHERE chat_id = $1
			ORDER BY id DESC
			LIMIT 1
		`, chat.ID).Scan(
			&msg.ID, &msg.Ciphertext, &msg.Nonce, &msg.ChatID, &msg.Seen, &msg.ReplyTo,
		)

		if err == nil {
			chat.LastMessage = &msg
		}

		chats = append(chats, &chat)
	}

	return chats, rows.Err()
}
