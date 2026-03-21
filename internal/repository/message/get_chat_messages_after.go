package message

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *MessageRepo) GetChatMessagesAfter(chatID, afterID, count int) ([]*domain.Message, error) {
	start := time.Now()

	rows, err := r.db.Query(`
	SELECT 
		id,
		ciphertext,
		nonce,
		chat_id,
		seen,
		reply_to
	FROM messages
	WHERE chat_id = $1 AND id > $2 
	ORDER BY id DESC
	LIMIT $3
	`, chatID, afterID, count)

	duration := time.Since(start)

	metrics.ObserveDB("messages_get_from_chat_after", duration, err)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(
			&message.ID,
			&message.Ciphertext,
			&message.Nonce,
			&message.ChatID,
			&message.Seen,
			&message.ReplyTo,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, rows.Err()
}
