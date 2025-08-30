package domain

import "time"

type Message struct {
	ID                    int        `db:"id" json:"id"`
	Ciphertext            string     `db:"ciphertext" json:"ciphertext"`
	EncapsulatedKey       string     `db:"encapsulated_key" json:"encapsulated_key"`
	Nonce                 string     `db:"nonce" json:"nonce"`
	ChatID                int        `db:"chat_id" json:"chat_id"`
	Signature             string     `db:"signature" json:"signature"`
	SignedPayload         string     `db:"signed_payload" json:"signed_payload"`
	CEKWrap               string     `db:"cek_wrap" json:"cek_wrap"`
	CEKWrapIV             string     `db:"cek_wrap_iv" json:"cek_wrap_iv"`
	CEKWrapSalt           string     `db:"cek_wrap_salt" json:"cek_wrap_salt"`
	EncapsulatedKeySender string     `db:"encapsulated_key_sender" json:"encapsulated_key_sender"`
	CEKWrapSender         string     `db:"cek_wrap_sender" json:"cek_wrap_sender"`
	CEKWrapSenderIV       string     `db:"cek_wrap_sender_iv" json:"cek_wrap_sender_iv"`
	CEKWrapSenderSalt     string     `db:"cek_wrap_sender_salt" json:"cek_wrap_sender_salt"`
	Seen                  *time.Time `db:"seen" json:"seen,omitempty"`
}

type SocketMessage struct {
	Type                  string `json:"type"`
	Ciphertext            string `json:"ciphertext"`
	EncapsulatedKey       string `json:"encapsulated_key"`
	Nonce                 string `json:"nonce"`
	ChatID                int    `json:"chat_id"`
	Signature             string `json:"signature"`
	SignedPayload         string `json:"signed_payload"`
	CEKWrap               string `json:"cek_wrap"`
	CEKWrapIV             string `json:"cek_wrap_iv"`
	CEKWrapSalt           string `json:"cek_wrap_salt"`
	EncapsulatedKeySender string `json:"encapsulated_key_sender"`
	CEKWrapSender         string `json:"cek_wrap_sender"`
	CEKWrapSenderIV       string `json:"cek_wrap_sender_iv"`
	CEKWrapSenderSalt     string `json:"cek_wrap_sender_salt"`
}
