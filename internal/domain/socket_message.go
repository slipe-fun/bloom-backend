package domain

type SocketMessage struct {
	ChatID         int    `json:"chat_id"`
	EncryptionType string `json:"encryption_type"`
	Ciphertext     string `json:"ciphertext"`
	Nonce          string `json:"nonce"`
	ReplyTo        int    `json:"reply_to,omitempty"`
}
