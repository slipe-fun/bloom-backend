package domain

import "time"

type User struct {
	ID             int       `db:"id" json:"-"`
	PublicID       string    `db:"public_id" json:"id"`
	Username       string    `db:"username" json:"username"`
	DisplayName    *string   `db:"display_name" json:"display_name"`
	Description    *string   `db:"description" json:"description"`
	MlKemPublicKey string    `db:"ml_kem_public_key" json:"ml_kem_public_key"`
	EcdhPublicKey  string    `db:"ecdh_public_key" json:"ecdh_public_key"`
	EdPublicKey    string    `db:"ed_public_key" json:"ed_public_key"`
	Date           time.Time `db:"date" json:"date"`
}
