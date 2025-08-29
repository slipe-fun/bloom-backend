package ChatRepo

import (
	"encoding/json"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *ChatRepo) GetById(id int) (*domain.Chat, error) {
	var chat domain.Chat
	var membersJSON []byte

	query := `SELECT id, members FROM chats WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&chat.ID, &membersJSON)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(membersJSON, &chat.Members)

	return &chat, nil
}
