package crypto

import (
	"encoding/base64"

	"github.com/cloudflare/circl/sign/ed448"
)

func VerifyEncryptedMasterKeySignature(
	pubKey []byte,
	signature []byte,
	ciphertext,
	nonce,
	salt,
	mlKem,
	ecdh,
	ed string,
) (bool, error) {

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return false, err
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return false, err
	}

	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false, err
	}

	mlKemBytes, err := base64.StdEncoding.DecodeString(mlKem)
	if err != nil {
		return false, err
	}

	ecdhBytes, err := base64.StdEncoding.DecodeString(ecdh)
	if err != nil {
		return false, err
	}

	edBytes, err := base64.StdEncoding.DecodeString(ed)
	if err != nil {
		return false, err
	}

	message := ConcatBytes(
		ciphertextBytes,
		nonceBytes,
		saltBytes,
		mlKemBytes,
		ecdhBytes,
		edBytes,
	)

	return ed448.Verify(pubKey, message, signature, ""), nil
}
