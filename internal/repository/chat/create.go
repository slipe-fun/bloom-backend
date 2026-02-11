package chat

import (
	"encoding/json"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *ChatRepo) Create(chat *domain.Chat) (*domain.Chat, error) {
	membersJSON, _ := json.Marshal(chat.Members)

	query := `INSERT INTO chats (members, encryption_key) VALUES ($1, $2) RETURNING id, members, encryption_key`

	var created domain.Chat
	var membersBytes []byte
	err := r.db.QueryRow(query, membersJSON, chat.EncryptionKey).Scan(&created.ID, &membersBytes, &created.EncryptionKey)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(membersBytes, &created.Members); err != nil {
		return nil, err
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
