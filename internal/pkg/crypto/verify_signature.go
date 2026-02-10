package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
)

func VerifySignature(publicKeyBase64 string, payload, signatureBase64 string) error {
	pubKey, err := DecodeEd25519Raw(publicKeyBase64)
	if err != nil {
		fmt.Print(err)
		return err
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		fmt.Print(err)
		return errors.New("не удалось декодировать подпись")
	}
	if len(sigBytes) != ed25519.SignatureSize {
		return errors.New("некорректная длина подписи")
	}

	if !ed25519.Verify(pubKey, []byte(payload), sigBytes) {
		return errors.New("подпись невалидна")
	}
	return nil
}
