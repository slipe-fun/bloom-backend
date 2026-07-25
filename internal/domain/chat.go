package domain

type EncryptedSyncKey struct {
	CipherText string `json:"ciphertext" db:"ciphertext"`
	Nonce      string `json:"nonce" db:"nonce"`
}

type Handshake struct {
	ReceiverCipherText  string `json:"receiver_cipher_text"`
	SenderEphemeralX448 string `json:"sender_ephemeral_x448"`
	EncryptedSyncKey    `json:"encrypted_sync_key"`
}

type Chat struct {
	ID        int        `json:"id"`
	Members   []User     `json:"members"`
	Handshake *Handshake `json:"handshake,omitempty"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
}

type RawChat struct {
	Members   []Member   `json:"members"`
	Handshake *Handshake `json:"handshake"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
}

type Member struct {
	ID int `json:"id"`
}

type ChatWithLastMessage struct {
	ID              int        `json:"id"`
	Members         []User     `json:"members"`
	Handshake       *Handshake `json:"handshake,omitempty"`
	LastMessage     *Message   `json:"last_message,omitempty"`
	LastReadMessage *Message   `json:"last_read_message,omitempty"`
}
