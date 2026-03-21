package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *MessageRepo) Create(message *domain.Message) (*domain.Message, error) {
	query := `INSERT INTO messages 
		(ciphertext, nonce, chat_id, reply_to) 
		VALUES ($1,$2,$3,$4) 
		RETURNING id, ciphertext, nonce, chat_id, reply_to`

	created := domain.Message{}

	err := r.db.QueryRow(
		query,
		message.Ciphertext,
		message.Nonce,
		message.ChatID,
		nullInt(message.ReplyTo),
	).Scan(
		&created.ID,
		&created.Ciphertext,
		&created.Nonce,
		&created.ChatID,
		&created.ReplyTo,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func nullInt(i *int) interface{} {
	if i == nil {
		return nil
	}
	return *i
}
