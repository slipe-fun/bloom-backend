package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
)

func DecodeEd25519Raw(pubKeyBase64 string) (ed25519.PublicKey, error) {
	raw, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		return nil, errors.New("не удалось декодировать публичный ключ")
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, errors.New("неверная длина публичного ключа")
	}
	return ed25519.PublicKey(raw), nil
}
