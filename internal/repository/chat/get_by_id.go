package chat

import (
	"encoding/json"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *ChatRepo) GetById(id int) (*domain.Chat, error) {
	var chat domain.Chat
	var membersJSON []byte

	query := `SELECT id, members, encryption_key FROM chats WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&chat.ID, &membersJSON, &chat.EncryptionKey)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(membersJSON, &chat.Members)

	for i := range chat.Members {
		member := chat.Members[i]
		user, err := r.userRepo.GetById(member.ID)
		if err != nil {
			continue
		}
		chat.Members[i].Username = user.Username
	}

	return &chat, nil
}
