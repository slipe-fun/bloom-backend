package domain

type EncryptedKeys struct {
	ID         int    `db:"id" json:"id"`
	UserID     int    `db:"user_id" json:"user_id"`
	Ciphertext string `db:"ciphertext" json:"ciphertext"`
	Nonce      string `db:"nonce" json:"nonce"`
	Salt       string `db:"salt" json:"salt"`
}
