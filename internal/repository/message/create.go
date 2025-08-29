package MessageRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

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
