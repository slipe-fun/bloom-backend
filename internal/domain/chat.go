package domain

import "time"

type Member struct {
	ID          int       `json:"id"`
	Username    string    `json:"username,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
}

type Chat struct {
	ID            int      `db:"id" json:"id"`
	Members       []Member `db:"members" json:"members"`
	EncryptionKey *string  `db:"encryption_key" json:"encryption_key"`
}

type ChatWithLastMessage struct {
	ID              int      `json:"id"`
	Members         []Member `json:"members"`
	EncryptionKey   *string  `json:"encryption_key"`
	LastMessage     *Message `json:"last_message,omitempty"`
	LastReadMessage *Message `json:"last_read_message,omitempty"`
}
