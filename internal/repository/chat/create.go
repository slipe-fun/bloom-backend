package ChatRepo

import (
	"encoding/json"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *ChatRepo) Create(chat *domain.Chat) (*domain.Chat, error) {
	membersJSON, _ := json.Marshal(chat.Members)

	query := `INSERT INTO chats (members) VALUES ($1) RETURNING id, members`

	var created domain.Chat
	var membersBytes []byte
	err := r.db.QueryRow(query, membersJSON).Scan(&created.ID, &membersBytes)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(membersBytes, &created.Members); err != nil {
		return nil, err
	}

	return &created, nil
}
