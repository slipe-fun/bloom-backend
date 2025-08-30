package domain

type SocketKeys struct {
	ChatID         int    `json:"chat_id"`
	KyberPublicKey string `json:"kyberPublicKey"`
	EcdhPublicKey  string `json:"ecdhPublicKey"`
	EdPublicKey    string `json:"edPublicKey"`
}
