package domain

import "time"

type Session struct {
	ID        int        `db:"id" json:"id"`
	Token     string     `db:"token" json:"token"`
	UserID    int        `db:"user_id" json:"-"`
	PublicID  string     `db:"user_public_id" json:"user_id"`
	RevokedAt *time.Time `db:"revoked_at" json:"revoked_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}
