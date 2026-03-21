package message

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *MessageRepo) GetByID(id int) (*domain.Message, error) {
	var message domain.Message

	query := `SELECT 
		id,
		ciphertext,
		nonce,
		chat_id,
		seen,
		reply_to
	FROM messages WHERE id = $1`

	start := time.Now()

	err := r.db.Get(&message, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("message_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &message, nil
}
