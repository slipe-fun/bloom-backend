package domain

import "time"

type EncryptedChatKeys struct {
	ID              int       `db:"id" json:"id"`
	ChatID          int       `db:"chat_id" json:"chat_id"`
	SessionID       int       `db:"session_id" json:"session_id"`
	EncryptedKey    string    `db:"encrypted_key" json:"encrypted_key"`
	EncapsulatedKey string    `db:"encapsulated_key" json:"encapsulated_key"`
	Nonce           string    `db:"nonce" json:"nonce"`
	Salt            string    `db:"salt" json:"salt"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type RawEncryptedChatKeys struct {
	SessionID       int    `json:"session_id"`
	EncryptedKey    string `json:"encrypted_key"`
	EncapsulatedKey string `json:"encapsulated_key"`
	Nonce           string `json:"nonce"`
	Salt            string `json:"salt"`
}
