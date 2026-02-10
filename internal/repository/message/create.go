package message

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (r *MessageRepo) Create(message *domain.Message) (*domain.Message, error) {
	query := `INSERT INTO messages 
		(ciphertext, nonce, chat_id, encapsulated_key, signature, signed_payload, cek_wrap, cek_wrap_iv, cek_wrap_salt, encapsulated_key_sender, cek_wrap_sender, cek_wrap_sender_iv, cek_wrap_sender_salt, reply_to) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) 
		RETURNING id, ciphertext, nonce, chat_id, encapsulated_key, signature, signed_payload, cek_wrap, cek_wrap_iv, cek_wrap_salt, encapsulated_key_sender, cek_wrap_sender, cek_wrap_sender_iv, cek_wrap_sender_salt, reply_to`

	created := domain.Message{}

	err := r.db.QueryRow(
		query,
		message.Ciphertext,
		message.Nonce,
		message.ChatID,
		nullStr(message.EncapsulatedKey),
		nullStr(message.Signature),
		nullStr(message.SignedPayload),
		nullStr(message.CEKWrap),
		nullStr(message.CEKWrapIV),
		nullStr(message.CEKWrapSalt),
		nullStr(message.EncapsulatedKeySender),
		nullStr(message.CEKWrapSender),
		nullStr(message.CEKWrapSenderIV),
		nullStr(message.CEKWrapSenderSalt),
		nullInt(message.ReplyTo),
	).Scan(
		&created.ID,
		&created.Ciphertext,
		&created.Nonce,
		&created.ChatID,
		&created.EncapsulatedKey,
		&created.Signature,
		&created.SignedPayload,
		&created.CEKWrap,
		&created.CEKWrapIV,
		&created.CEKWrapSalt,
		&created.EncapsulatedKeySender,
		&created.CEKWrapSender,
		&created.CEKWrapSenderIV,
		&created.CEKWrapSenderSalt,
		&created.ReplyTo,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func nullStr(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func nullInt(i *int) interface{} {
	if i == nil {
		return nil
	}
	return *i
}
