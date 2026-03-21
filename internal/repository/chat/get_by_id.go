package chat

import (
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) GetByID(id int) (*domain.Chat, error) {
	var chat domain.Chat
	var membersJSON []byte

	query := `SELECT id, members, encryption_key FROM chats WHERE id=$1`

	start := time.Now()

	err := r.db.QueryRow(query, id).Scan(&chat.ID, &membersJSON, &chat.EncryptionKey)

	duration := time.Since(start)

	metrics.ObserveDB("chat_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(membersJSON, &chat.Members)

	for i := range chat.Members {
		member := chat.Members[i]
		user, err := r.userRepo.GetByID(member.ID)
		if err != nil {
			continue
		}
		chat.Members[i].Username = user.Username

		if user.DisplayName != nil {
			chat.Members[i].DisplayName = *user.DisplayName
		}

		if user.Description != nil {
			chat.Members[i].Description = *user.Description
		}

		chat.Members[i].Date = user.Date
	}

	return &chat, nil
}
