package crypto

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

func DecodeEd25519SPKI(pubKeyBase64 string) (ed25519.PublicKey, error) {
	der, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		return nil, errors.New("не удалось декодировать публичный ключ")
	}

	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, errors.New("не удалось распарсить SPKI публичный ключ")
	}

	edKey, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("не Ed25519 публичный ключ")
	}

	return edKey, nil
}
