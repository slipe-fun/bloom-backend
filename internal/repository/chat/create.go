package chat

import (
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *ChatRepo) Create(chat *domain.RawChat) (*domain.Chat, error) {
	var handshakeJSON any = nil
	var membersJSON any = nil

	if chat.Type == "private" {
		mBytes, _ := json.Marshal(chat.Members)
		membersJSON = mBytes

		if chat.Handshake != nil {
			hBytes, _ := json.Marshal(chat.Handshake)
			handshakeJSON = hBytes
		}
	}

	query := `INSERT INTO chats (members, handshake, type, title) VALUES ($1, $2, $3, $4) RETURNING id, members, handshake, type, title`

	var created domain.Chat
	var membersBytes []byte
	var handshakeBytes []byte
	var titlePtr *string

	start := time.Now()

	err := r.db.QueryRow(query, membersJSON, handshakeJSON, chat.Type, chat.Title).
		Scan(&created.ID, &membersBytes, &handshakeBytes, &created.Type, &titlePtr)

	duration := time.Since(start)

	metrics.ObserveDB("chat_create", duration, err)

	if err != nil {
		return nil, err
	}

	if titlePtr != nil {
		created.Title = *titlePtr
	}

	if chat.Type == "private" {
		var rawMembers []domain.Member
		if len(membersBytes) > 0 {
			if err := json.Unmarshal(membersBytes, &rawMembers); err != nil {
				return nil, err
			}
		}

		created.Members = make([]domain.User, 0, len(rawMembers))

		if len(handshakeBytes) > 0 {
			var hs domain.Handshake
			if err := json.Unmarshal(handshakeBytes, &hs); err == nil {
				created.Handshake = &hs
			}
		}

		for _, m := range rawMembers {
			user, err := r.userRepo.GetByID(m.ID)
			if err != nil {
				continue
			}

			created.Members = append(created.Members, *user)
		}
	}

	return &created, nil
}
