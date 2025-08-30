package MessageRepo

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *MessageRepo) GetChatMessages(id int) ([]*domain.Message, error) {
	rows, err := r.db.Query(`
	SELECT 
		id,
		ciphertext,
		encapsulated_key,
		nonce,
		chat_id,
		signature,
		seen,
		COALESCE(signed_payload, '') AS signed_payload,
		COALESCE(cek_wrap, '') AS cek_wrap,
		COALESCE(cek_wrap_iv, '') AS cek_wrap_iv,
		COALESCE(cek_wrap_salt, '') AS cek_wrap_salt,
		COALESCE(encapsulated_key_sender, '') AS encapsulated_key_sender,
		COALESCE(cek_wrap_sender, '') AS cek_wrap_sender,
		COALESCE(cek_wrap_sender_iv, '') AS cek_wrap_sender_iv,
		COALESCE(cek_wrap_sender_salt, '') AS cek_wrap_sender_salt
	FROM messages
	WHERE chat_id = $1
	`, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(
			&message.ID,
			&message.Ciphertext,
			&message.EncapsulatedKey,
			&message.Nonce,
			&message.ChatID,
			&message.Signature,
			&message.Seen,
			&message.SignedPayload,
			&message.CEKWrap,
			&message.CEKWrapIV,
			&message.CEKWrapSalt,
			&message.EncapsulatedKeySender,
			&message.CEKWrapSender,
			&message.CEKWrapSenderIV,
			&message.CEKWrapSenderSalt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, rows.Err()
}
