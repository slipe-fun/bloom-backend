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
		SELECT id, members, encryption_key
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
		if err := rows.Scan(&chat.ID, &membersJSON, &chat.EncryptionKey); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(membersJSON, &chat.Members); err != nil {
			return nil, err
		}

		for i := range chat.Members {
			member := chat.Members[i]
			user, err := r.userRepo.GetByID(member.ID)
			if err != nil {
				continue
			}
			chat.Members[i].Username = user.Username
		}

		var msg domain.Message
		err = r.db.QueryRow(`
			SELECT id, ciphertext, nonce, chat_id, encapsulated_key, signature,
			       signed_payload, cek_wrap, cek_wrap_iv, cek_wrap_salt,
				   encapsulated_key_sender, cek_wrap_sender, cek_wrap_sender_iv,
				   cek_wrap_sender_salt, seen, reply_to
			FROM messages
			WHERE chat_id = $1
			ORDER BY id DESC
			LIMIT 1
		`, chat.ID).Scan(
			&msg.ID, &msg.Ciphertext, &msg.Nonce, &msg.ChatID, &msg.EncapsulatedKey,
			&msg.Signature, &msg.SignedPayload, &msg.CEKWrap, &msg.CEKWrapIV, &msg.CEKWrapSalt,
			&msg.EncapsulatedKeySender, &msg.CEKWrapSender, &msg.CEKWrapSenderIV, &msg.CEKWrapSenderSalt,
			&msg.Seen, &msg.ReplyTo,
		)
		if err == nil {
			chat.LastMessage = &msg
		}

		chats = append(chats, &chat)
	}

	return chats, rows.Err()
}
