package domain

type UploadIdentityRequest struct {
	Ciphertext     string `json:"ciphertext"`
	Nonce          string `json:"nonce"`
	Signature      string `json:"signature"`
	MlKemPublicKey string `json:"ml_kem_public_key"`
	EcdhPublicKey  string `json:"ecdh_public_key"`
	EdPublicKey    string `json:"ed_public_key"`
	Salt           string `json:"salt"`
}
