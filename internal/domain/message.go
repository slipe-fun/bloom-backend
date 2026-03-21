package domain

import "time"

type Message struct {
	ID         int        `db:"id" json:"id"`
	Ciphertext string     `db:"ciphertext" json:"ciphertext"`
	Nonce      string     `db:"nonce" json:"nonce"`
	ChatID     int        `db:"chat_id" json:"chat_id"`
	Seen       *time.Time `db:"seen" json:"seen,omitempty"`
	ReplyTo    *int       `db:"reply_to" json:"reply_to,omitempty"`
}

type MessageWithReply struct {
	Message
	ReplyToMessage *Message `json:"reply_to,omitempty"`
}
