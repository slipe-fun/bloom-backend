package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type MessageRepo struct {
	db *sqlx.DB
}

func NewMessageRepo(db *sqlx.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Create(message *domain.Message) error {
	query := `INSERT INTO messages 
		(ciphertext, encapsulated_key, nonce, chat_id, signature, salt) VALUES 
		($1, $2, $3, $4, $5, $6) RETURNING id`
	return r.db.QueryRow(query,
		message.Ciphertext,
		message.EncapsulatedKey,
		message.Nonce,
		message.ChatID,
		message.Signature,
		message.Salt,
	).Scan(&message.ID)
}

func (r *MessageRepo) GetByChatID(id int) (*domain.Message, error) {
	var message domain.Message

	query := `SELECT id, ciphertext, encapsulated_key, nonce, chat_id, signature, salt FROM messages WHERE id = $1`
	err := r.db.Get(&message, query, id)

	if err != nil {
		return nil, err
	}

	return &message, nil
}
