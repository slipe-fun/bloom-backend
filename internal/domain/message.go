package domain

type Message struct {
	ID              int    `db:"id"`
	Ciphertext      string `db:"ciphertext"`
	EncapsulatedKey string `db:"encapsulated_key"`
	Nonce           string `db:"nonche"`
	ChatID          int    `db:"chat_id"`
	Signature       string `db:"signature"`
	Salt            string `db:"salt"`
}
