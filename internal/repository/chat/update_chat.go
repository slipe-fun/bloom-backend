package chat

import (
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) UpdateChat(chat *domain.Chat) error {
	membersJSON, err := json.Marshal(chat.Members)
	if err != nil {
		return err
	}

	start := time.Now()

	_, err = r.db.Exec(`
        UPDATE chats
        SET members = $1
        WHERE id = $2
    `, membersJSON, chat.ID)

	duration := time.Since(start)

	metrics.ObserveDB("chat_update", duration, err)

	if err != nil {
		return err
	}

	return nil
}
