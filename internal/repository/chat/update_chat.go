package chat

import (
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) UpdateChat(chat *domain.Chat) error {
	rawMembers := make([]domain.Member, len(chat.Members))

	for i, m := range chat.Members {
		rawMembers[i] = domain.Member{
			ID: m.ID,
		}
	}

	membersJSON, err := json.Marshal(rawMembers)
	if err != nil {
		return err
	}

	var handshakeJSON []byte
	if chat.Handshake != nil {
		var err error
		handshakeJSON, err = json.Marshal(chat.Handshake)
		if err != nil {
			return err
		}
	}

	start := time.Now()

	_, err = r.db.Exec(`
		UPDATE chats
		SET members = $1, handshake = $2
		WHERE id = $3
	`, membersJSON, handshakeJSON, chat.ID)

	duration := time.Since(start)

	metrics.ObserveDB("chat_update", duration, err)

	return err
}
