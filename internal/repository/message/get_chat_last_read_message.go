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
		nonce,
		chat_id,
		seen,
		reply_to
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
