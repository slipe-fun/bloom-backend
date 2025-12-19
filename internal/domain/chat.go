package domain

type Member struct {
	ID             int    `json:"id"`
	Username       string `json:"username,omitempty"`
	KyberPublicKey string `json:"kyberPublicKey"`
	EcdhPublicKey  string `json:"ecdhPublicKey"`
	EdPublicKey    string `json:"edPublicKey"`
}

type Chat struct {
	ID            int      `db:"id" json:"id"`
	Members       []Member `db:"members" json:"members"`
	EncryptionKey *string  `db:"encryption_key" json:"encryption_key"`
}

type ChatWithLastMessage struct {
	ID            int      `json:"id"`
	Members       []Member `json:"members"`
	EncryptionKey *string  `json:"encryption_key"`
	LastMessage   *Message `json:"last_message,omitempty"`
}
