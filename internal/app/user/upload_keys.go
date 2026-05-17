package user

import (
	"encoding/base64"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (u *UserApp) UploadIdentityKeys(user_id int, req *domain.UploadIdentityRequest) error {
	user, err := u.users.GetByID(user_id)
	if err != nil || user == nil {
		return domain.Failed("user not found")
	}

	if user.KyberPublicKey != "" || user.EcdhPublicKey != "" || user.EdPublicKey != "" {
		return domain.InvalidData("user already have keys")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil || len(sigBytes) != 114 {
		return domain.InvalidData("invalid signature length")
	}

	verifyKeyBytes, err := base64.StdEncoding.DecodeString(req.EdPublicKey)
	if err != nil || len(verifyKeyBytes) != 57 {
		return domain.InvalidData("invalid public key length")
	}

	ecdhBytes, err := base64.StdEncoding.DecodeString(req.EcdhPublicKey)
	if err != nil || len(ecdhBytes) != 56 {
		return domain.InvalidData("invalid ecdh public key")
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(req.Nonce)
	if err != nil || len(nonceBytes) != 12 {
		return domain.InvalidData("invalid nonce")
	}

	saltBytes, err := base64.StdEncoding.DecodeString(req.Salt)
	if err != nil || len(saltBytes) != 16 {
		return domain.InvalidData("invalid salt")
	}

	mlKemBytes, err := base64.StdEncoding.DecodeString(req.MlKemPublicKey)
	if err != nil || len(mlKemBytes) != 1184 {
		return domain.InvalidData("invalid ml-kem public key")
	}

	if len(req.Ciphertext) > 10000 {
		return domain.InvalidData("ciphertext too large")
	}
	_, err = base64.StdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		return domain.InvalidData("invalid ciphertext format")
	}

	isValid, err := crypto.VerifyIdentityKeysSignature(
		verifyKeyBytes,
		sigBytes,
		req.Ciphertext,
		req.Nonce,
		req.MlKemPublicKey,
		req.EcdhPublicKey,
		req.EdPublicKey,
		req.Salt,
	)
	if err != nil {
		return domain.Failed("error during verification")
	}
	if !isValid {
		return domain.InvalidData("invalid identity signature")
	}

	_, err = u.keys.Create(&domain.EncryptedKeys{
		UserID:     user_id,
		Type:       "bundle",
		Ciphertext: req.Ciphertext,
		Nonce:      req.Nonce,
		Signature:  req.Signature,
		Salt:       req.Salt,
	})
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return domain.Failed("failed to save encrypted bundle")
	}

	err = u.users.UpdatePublicKeys(user_id, req.MlKemPublicKey, req.EcdhPublicKey, req.EdPublicKey)
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return domain.Failed("failed to update user public keys")
	}

	return nil
}
