package chat

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) GetByIDs(ids []int) ([]*domain.Chat, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `SELECT id, members, handshake FROM chats WHERE id = ANY($1)`

	start := time.Now()
	rows, err := r.db.Query(query, pq.Array(ids))
	duration := time.Since(start)
	metrics.ObserveDB("chat_get_by_ids", duration, err)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*domain.Chat
	for rows.Next() {
		var chat domain.Chat
		var membersJSON []byte
		var handshakeJSON []byte

		if err := rows.Scan(&chat.ID, &membersJSON, &handshakeJSON); err != nil {
			return nil, err
		}

		json.Unmarshal(membersJSON, &chat.Members)

		if len(handshakeJSON) > 0 {
			var hs domain.Handshake
			if err := json.Unmarshal(handshakeJSON, &hs); err == nil {
				chat.Handshake = &hs
			}
		}

		chats = append(chats, &chat)
	}

	return chats, nil
}
