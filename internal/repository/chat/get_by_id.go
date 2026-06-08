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
	var handshakeJSON []byte

	query := `SELECT id, members, handshake FROM chats WHERE id=$1`

	start := time.Now()

	err := r.db.QueryRow(query, id).Scan(&chat.ID, &membersJSON, &handshakeJSON)

	duration := time.Since(start)

	metrics.ObserveDB("chat_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	var rawMembers []domain.Member
	if err := json.Unmarshal(membersJSON, &rawMembers); err != nil {
		return nil, err
	}

	chat.Members = make([]domain.User, 0, len(rawMembers))

	for _, m := range rawMembers {
		user, err := r.userRepo.GetByID(m.ID)
		if err != nil {
			continue
		}
		chat.Members = append(chat.Members, *user)
	}

	if len(handshakeJSON) > 0 {
		var hs domain.Handshake
		if err := json.Unmarshal(handshakeJSON, &hs); err == nil {
			chat.Handshake = &hs
		}
	}

	return &chat, nil
}
