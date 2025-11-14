package domain

import "time"

type VerificationCode struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Code      string    `db:"code" json:"code"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}
