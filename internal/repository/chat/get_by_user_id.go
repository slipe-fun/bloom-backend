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
		SELECT id, members, handshake
		FROM chats
		WHERE EXISTS (
			SELECT 1 FROM jsonb_array_elements(members) AS m
			WHERE (m->>'id')::int = $1
		)
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

		if err := rows.Scan(&chat.ID, &membersJSON, &handshakeJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(membersJSON, &chat.Members); err != nil {
			return nil, err
		}

		if len(handshakeJSON) > 0 {
			var hs domain.Handshake
			if err := json.Unmarshal(handshakeJSON, &hs); err == nil {
				chat.Handshake = &hs
			}
		}

		for i := range chat.Members {
			member := chat.Members[i]
			user, err := r.userRepo.GetByID(member.ID)
			if err != nil {
				continue
			}
			chat.Members[i] = *user
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
