package service

import (
	"encoding/base64"
	"errors"
)

func CheckKeysLength(kyberPublicKey string, ecdhPublicKey string, edPublicKey string) error {
	kyberKey, err := base64.StdEncoding.DecodeString(kyberPublicKey)
	if err != nil {
		return errors.New("invalid base64 for Kyber key")
	}
	ecdhKey, err := base64.StdEncoding.DecodeString(ecdhPublicKey)
	if err != nil {
		return errors.New("invalid base64 for ECDH key")
	}
	edKey, err := base64.StdEncoding.DecodeString(edPublicKey)
	if err != nil {
		return errors.New("invalid base64 for Ed25519 key")
	}

	if len(kyberKey) != 1184 {
		return errors.New("invalid Kyber key length")
	}
	if len(ecdhKey) != 65 {
		return errors.New("invalid ECDH key length")
	}
	if len(edKey) != 44 {
		return errors.New("invalid Ed25519 key length")
	}

	return nil
}
