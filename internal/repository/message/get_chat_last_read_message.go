package message

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *MessageRepo) GetChatLastReadMessage(chatID int) (*domain.Message, error) {
	var message domain.Message

	query := `SELECT 
		id,
		ciphertext,
		encapsulated_key,
		nonce,
		chat_id,
		signature,
		seen,
		COALESCE(signed_payload, '') AS signed_payload,
		COALESCE(cek_wrap, '') AS cek_wrap,
		COALESCE(cek_wrap_iv, '') AS cek_wrap_iv,
		COALESCE(cek_wrap_salt, '') AS cek_wrap_salt,
		COALESCE(encapsulated_key_sender, '') AS encapsulated_key_sender,
		COALESCE(cek_wrap_sender, '') AS cek_wrap_sender,
		COALESCE(cek_wrap_sender_iv, '') AS cek_wrap_sender_iv,
		COALESCE(cek_wrap_sender_salt, '') AS cek_wrap_sender_salt,
		COALESCE(reply_to, 0) AS reply_to
	FROM messages WHERE chat_id = $1 AND seen IS NOT NULL
	ORDER BY id DESC
	LIMIT 1`

	start := time.Now()

	err := r.db.Get(&message, query, chatID)

	duration := time.Since(start)

	metrics.ObserveDB("message_get_chat_last_read", duration, err)

	if err != nil {
		return nil, err
	}

	return &message, nil
}
