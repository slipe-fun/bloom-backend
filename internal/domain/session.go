package domain

import "time"

type Session struct {
	ID                int        `db:"id" json:"id"`
	Token             string     `db:"token" json:"token"`
	UserID            int        `db:"user_id" json:"user_id"`
	IdentityPublicKey *string    `db:"identity_pub" json:"identity_pub"`
	EcdhPublicKey     *string    `db:"ecdh_pub" json:"ecdh_pub"`
	KyberPublicKey    *string    `db:"kyber_pub" json:"kyber_pub"`
	RevokedAt         *time.Time `db:"revoked_at" json:"revoked_at"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
}
