package ChatRepo

import (
	"encoding/json"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *ChatRepo) GetByUserId(id int) ([]*domain.Chat, error) {
	rows, err := r.db.Query(`
		SELECT id, members
		FROM chats
		WHERE EXISTS (
			SELECT 1 FROM jsonb_array_elements(members) AS m
			WHERE (m->>'id')::int = $1
		)
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*domain.Chat
	for rows.Next() {
		var chat domain.Chat
		var membersJSON []byte
		if err := rows.Scan(&chat.ID, &membersJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(membersJSON, &chat.Members); err != nil {
			return nil, err
		}
		chats = append(chats, &chat)
	}
	return chats, rows.Err()
}
