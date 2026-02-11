package chat

import (
	"encoding/json"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *ChatRepo) GetWithUsers(id int, recipient int) (*domain.Chat, error) {
	var chat domain.Chat
	var membersJSON []byte

	query := `
	SELECT id, members, encryption_key
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
	err := r.db.QueryRow(query, id, recipient).Scan(&chat.ID, &membersJSON, &chat.EncryptionKey)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(membersJSON, &chat.Members)

	for i := range chat.Members {
		member := chat.Members[i]
		user, err := r.userRepo.GetByID(member.ID)
		if err != nil {
			continue
		}
		chat.Members[i].Username = user.Username
	}

	return &chat, nil
}
