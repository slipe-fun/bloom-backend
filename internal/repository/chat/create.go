package chat

import (
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) Create(chat *domain.Chat) (*domain.Chat, error) {
	membersJSON, _ := json.Marshal(chat.Members)

	var handshakeJSON []byte
	if chat.Handshake != nil {
		handshakeJSON, _ = json.Marshal(chat.Handshake)
	}

	query := `INSERT INTO chats (members, handshake) VALUES ($1, $2) RETURNING id, members, handshake`

	var created domain.Chat
	var membersBytes []byte
	var handshakeBytes []byte

	start := time.Now()

	err := r.db.QueryRow(query, membersJSON, handshakeJSON).Scan(&created.ID, &membersBytes, &handshakeBytes)

	duration := time.Since(start)

	metrics.ObserveDB("chat_create", duration, err)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(membersBytes, &created.Members); err != nil {
		return nil, err
	}

	if len(handshakeBytes) > 0 {
		var hs domain.Handshake
		if err := json.Unmarshal(handshakeBytes, &hs); err == nil {
			created.Handshake = &hs
		}
	}

	for i := range created.Members {
		member := created.Members[i]
		user, err := r.userRepo.GetByID(member.ID)
		if err != nil {
			continue
		}
		created.Members[i].Username = user.Username
	}

	return &created, nil
}
