package crypto

import (
	"encoding/json"

	"github.com/cloudflare/circl/sign/ed448"
)

type IdentityKeysPayload struct {
	Ciphertext     string `json:"ciphertext"`
	Nonce          string `json:"nonce"`
	MlKemPublicKey string `json:"ml_kem_public_key"`
	EcdhPublicKey  string `json:"ecdh_public_key"`
	EdPublicKey    string `json:"ed_public_key"`
	Salt           string `json:"salt"`
}

func VerifyIdentityKeysSignature(pubKey []byte, signature []byte, ciphertext, nonce, mlKem, ecdh, ed, salt string) (bool, error) {
	data := IdentityKeysPayload{
		Ciphertext:     ciphertext,
		Nonce:          nonce,
		MlKemPublicKey: mlKem,
		EcdhPublicKey:  ecdh,
		EdPublicKey:    ed,
		Salt:           salt,
	}

	message, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	isValid := ed448.Verify(pubKey, message, signature, "")
	return isValid, nil
}
