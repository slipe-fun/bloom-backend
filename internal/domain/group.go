package domain

import "time"

type EncryptedGroupKey struct {
	Ciphertext string `json:"ciphertext" db:"ciphertext"`
	Nonce      string `json:"nonce" db:"nonce"`
	Salt       string `json:"salt" db:"salt"`
}

type GroupMember struct {
	ChatID            int       `json:"chat_id" db:"chat_id"`
	MemberID          int       `json:"member_id" db:"member_id"`
	InvitedByID       int       `json:"invited_by_id" db:"invited_by_id"`
	Role              string    `json:"role" db:"role"`
	JoinedAt          time.Time `json:"joined_at" db:"joined_at"`
	Handshake         `json:"handshake"`
	EncryptedGroupKey `json:"encrypted_group_key"`
}
