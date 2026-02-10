package chat

import (
	"encoding/json"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *ChatRepo) UpdateChat(chat *domain.Chat) error {
	membersJSON, err := json.Marshal(chat.Members)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
        UPDATE chats
        SET members = $1
        WHERE id = $2
    `, membersJSON, chat.ID)

	if err != nil {
		return err
	}

	return nil
}
