package domain

type Handshake struct {
	ReceiverCipherText string `json:"receiver_cipher_text"`
	SenderCipherText   string `json:"sender_cipher_text"`
	EncryptedSyncKey   struct {
		CipherText string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
	} `json:"encrypted_sync_key"`
}

type Chat struct {
	ID        int        `json:"id"`
	Members   []User     `json:"members"`
	Handshake *Handshake `json:"handshake,omitempty"`
}
type ChatWithLastMessage struct {
	ID              int        `json:"id"`
	Members         []User     `json:"members"`
	Handshake       *Handshake `json:"handshake,omitempty"`
	LastMessage     *Message   `json:"last_message,omitempty"`
	LastReadMessage *Message   `json:"last_read_message,omitempty"`
}
