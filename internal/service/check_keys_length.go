package service

import (
	"encoding/base64"
	"errors"
)

func CheckKeysLength(kyberPublicKey string, ecdhPublicKey string, edPublicKey string) error {
	kyberKey, err := base64.StdEncoding.DecodeString(kyberPublicKey)
	if err != nil {
		return errors.New("invalid_kyber_base64")
	}
	ecdhKey, err := base64.StdEncoding.DecodeString(ecdhPublicKey)
	if err != nil {
		return errors.New("invalid_ecdh_base64")
	}
	edKey, err := base64.StdEncoding.DecodeString(edPublicKey)
	if err != nil {
		return errors.New("invalid_ed25519_base64")
	}

	if len(kyberKey) != 1184 {
		return errors.New("invalid_kyber_key_length")
	}
	if len(ecdhKey) != 44 {
		return errors.New("invalid_ecdh_key_length")
	}
	if len(edKey) != 44 {
		return errors.New("invalid_ed25519_key_length")
	}

	return nil
}
