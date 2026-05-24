package chat

import (
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) GetWithUsers(id int, recipient int) (*domain.Chat, error) {
	var chat domain.Chat
	var membersJSON []byte
	var handshakeJSON []byte

	query := `
	SELECT id, members, handshake
	FROM chats
	WHERE EXISTS (
		SELECT 1
		FROM jsonb_array_elements(members) AS m
		WHERE (m->>'id')::int = $1
	)
	AND EXISTS (
		SELECT 1
		FROM jsonb_array_elements(members) AS m
		WHERE (m->>'id')::int = $2
	);
	`

	start := time.Now()

	err := r.db.QueryRow(query, id, recipient).Scan(&chat.ID, &membersJSON, &handshakeJSON)

	duration := time.Since(start)

	metrics.ObserveDB("chat_get_with_bot_users", duration, err)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(membersJSON, &chat.Members)

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

	return &chat, nil
}
