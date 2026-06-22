package domain

type Handshake struct {
	ReceiverCipherText  string `json:"receiver_cipher_text"`
	SenderEphemeralX448 string `json:"sender_ephemeral_x448"`
	EncryptedSyncKey    struct {
		CipherText string `json:"ciphertext"`
		Nonce      string `json:"nonce"`
	} `json:"encrypted_sync_key"`
}

type Chat struct {
	ID        int        `json:"id"`
	Members   []User     `json:"members"`
	Handshake *Handshake `json:"handshake,omitempty"`
}

type RawChat struct {
	Members   []Member   `json:"members"`
	Handshake *Handshake `json:"handshake"`
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
