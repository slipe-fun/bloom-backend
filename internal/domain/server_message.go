package domain

import "time"

type ServerMessage struct {
	ID        int       `json:"id" db:"id"`
	ServerID  int       `json:"server_id" db:"server_id"`
	MemberID  int       `json:"member_id" db:"member_id"`
	ChannelID int       `json:"channel_id" db:"channel_id"`
	Content   string    `json:"content" db:"content"`
	SentAt    time.Time `json:"sent_at" db:"sent_at"`
}
