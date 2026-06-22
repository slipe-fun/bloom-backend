package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *MessageRepo) Create(message *domain.Message) (*domain.Message, error) {
	query := `INSERT INTO messages
		(ciphertext, nonce, salt, chat_id, reply_to)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, ciphertext, nonce, salt, chat_id, reply_to`

	created := domain.Message{}

	err := r.db.QueryRow(
		query,
		message.Ciphertext,
		message.Nonce,
		message.Salt,
		message.ChatID,
		nullInt(message.ReplyTo),
	).Scan(
		&created.ID,
		&created.Ciphertext,
		&created.Nonce,
		&created.Salt,
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
