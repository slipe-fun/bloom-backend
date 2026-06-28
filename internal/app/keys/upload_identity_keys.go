package keys

import (
	"encoding/base64"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *KeysApp) UploadIdentityKeys(user_id int, req *domain.IdentityKeysRequest) (*domain.IdentityPublicKeysBytes, *domain.EncryptedKeyBytes, error) {
	user, err := k.users.GetByID(user_id)
	if err != nil || user == nil {
		return nil, nil, domain.Failed("user not found")
	}

	if user.MlKemPublicKey != "" || user.EcdhPublicKey != "" || user.EdPublicKey != "" {
		return nil, nil, domain.InvalidData("user already have keys")
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(req.EncryptedSecretKeys.Signature)
	if err != nil || len(signatureBytes) != 114 {
		return nil, nil, domain.InvalidData("invalid signature length")
	}

	edKeyBytes, err := base64.StdEncoding.DecodeString(req.IdentityPublicKeys.EdPublicKey)
	if err != nil || len(edKeyBytes) != 57 {
		return nil, nil, domain.InvalidData("invalid public key length")
	}

	ecdhBytes, err := base64.StdEncoding.DecodeString(req.IdentityPublicKeys.EcdhPublicKey)
	if err != nil || len(ecdhBytes) != 56 {
		return nil, nil, domain.InvalidData("invalid ecdh public key")
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(req.EncryptedSecretKeys.Nonce)
	if err != nil || len(nonceBytes) != 12 {
		return nil, nil, domain.InvalidData("invalid nonce")
	}

	saltBytes, err := base64.StdEncoding.DecodeString(req.EncryptedSecretKeys.Salt)
	if err != nil || len(saltBytes) != 32 {
		return nil, nil, domain.InvalidData("invalid salt")
	}

	mlKemBytes, err := base64.StdEncoding.DecodeString(req.IdentityPublicKeys.MlKemPublicKey)
	if err != nil || len(mlKemBytes) != 1184 {
		return nil, nil, domain.InvalidData("invalid ml-kem public key")
	}

	if len(req.EncryptedSecretKeys.Ciphertext) > 10000 {
		return nil, nil, domain.InvalidData("ciphertext too large")
	}
	ciphertextBytes, err := base64.StdEncoding.DecodeString(req.EncryptedSecretKeys.Ciphertext)
	if err != nil {
		return nil, nil, domain.InvalidData("invalid ciphertext format")
	}

	isValid, err := crypto.VerifyIdentityKeysSignature(
		edKeyBytes,
		signatureBytes,
		req.EncryptedSecretKeys.Ciphertext,
		req.EncryptedSecretKeys.Nonce,
		req.IdentityPublicKeys.MlKemPublicKey,
		req.IdentityPublicKeys.EcdhPublicKey,
		req.IdentityPublicKeys.EdPublicKey,
		req.EncryptedSecretKeys.Salt,
	)
	if err != nil {
		return nil, nil, domain.Failed("error during verification")
	}
	if !isValid {
		return nil, nil, domain.InvalidData("invalid identity signature")
	}

	_, err = k.keys.Create(&domain.EncryptedKeys{
		UserID: user_id,
		Type:   "bundle",
		EncryptedKey: domain.EncryptedKey{
			Ciphertext: req.EncryptedSecretKeys.Ciphertext,
			Nonce:      req.EncryptedSecretKeys.Nonce,
			Signature:  req.EncryptedSecretKeys.Signature,
			Salt:       req.EncryptedSecretKeys.Salt,
		},
	})
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, nil, domain.Failed("failed to save encrypted bundle")
	}

	err = k.users.UpdatePublicKeys(user_id, req.IdentityPublicKeys.MlKemPublicKey, req.IdentityPublicKeys.EcdhPublicKey, req.IdentityPublicKeys.EdPublicKey)
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, nil, domain.Failed("failed to update user public keys")
	}

	return &domain.IdentityPublicKeysBytes{
			MlKemPublicKey: mlKemBytes,
			EcdhPublicKey:  ecdhBytes,
			EdPublicKey:    edKeyBytes,
		}, &domain.EncryptedKeyBytes{
			Ciphertext: ciphertextBytes,
			Nonce:      nonceBytes,
			Salt:       saltBytes,
			Signature:  signatureBytes,
		}, nil
}
