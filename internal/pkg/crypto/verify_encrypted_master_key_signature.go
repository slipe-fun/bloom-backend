package crypto

import (
	"encoding/json"
	"fmt"

	"github.com/cloudflare/circl/sign/ed448"
)

type EncryptedMasterKey struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Salt       string `json:"salt"`
}

func VerifyEncryptedMasterKeySignature(pubKey []byte, signature []byte, ciphertext, nonce, salt string) (bool, error) {
	data := EncryptedMasterKey{
		Ciphertext: ciphertext,
		Nonce:      nonce,
		Salt:       salt,
	}

	message, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	fmt.Printf("GO JSON STRING: %s\n", message)

	isValid := ed448.Verify(pubKey, message, signature, "")

	return isValid, nil
}
