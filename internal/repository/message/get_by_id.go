package MessageRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (r *MessageRepo) GetById(id int) (*domain.Message, error) {
	var message domain.Message

	query := `SELECT id, ciphertext, encapsulated_key, nonce, chat_id, signature, salt FROM messages WHERE id = $1`
	err := r.db.Get(&message, query, id)

	if err != nil {
		return nil, err
	}

	return &message, nil
}
