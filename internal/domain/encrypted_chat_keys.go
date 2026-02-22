package domain

import "time"

type EncryptedChatKeys struct {
	ID              int       `db:"id" json:"id"`
	ChatID          int       `db:"chat_id" json:"chat_id"`
	SessionID       int       `db:"session_id" json:"session_id"`
	FromSessionID   int       `db:"from_session_id" json:"from_session_id"`
	EncryptedKey    string    `db:"encrypted_key" json:"encrypted_key"`
	EncapsulatedKey string    `db:"encapsulated_key" json:"encapsulated_key"`
	CekWrap         string    `db:"cek_wrap" json:"cek_wrap"`
	CekWrapIV       string    `db:"cek_wrap_iv" json:"cek_wrap_iv"`
	Salt            string    `db:"salt" json:"salt"`
	Nonce           string    `db:"nonce" json:"nonce"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type RawEncryptedChatKeys struct {
	SessionID       int    `json:"session_id"`
	EncryptedKey    string `json:"encrypted_key"`
	EncapsulatedKey string `json:"encapsulated_key"`
	CekWrap         string `json:"cek_wrap"`
	CekWrapIV       string `json:"cek_wrap_iv"`
	Salt            string `json:"salt"`
	Nonce           string `json:"nonce"`
}
