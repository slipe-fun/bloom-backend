package auth

import (
	"context"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/generator"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (a *AuthApp) LoginBegin(authLookupID string) (*domain.KeysRequest, string, string, error) {
	user, err := a.users.GetByAuthLookupID(authLookupID)
	if err != nil {
		return nil, "", "", domain.NotFound("user is not found")
	}

	ctx := context.Background()

	challenge, err := generator.GenerateChallenge()
	if err != nil {
		return nil, "", "", domain.NotFound("failed to generate challenge")
	}

	err = a.rdb.Set(ctx, "auth:challenge:"+user.PublicID, challenge, 2*time.Minute).Err()
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return nil, "", "", domain.Failed("failed to save challenge")
	}

	encrypted_identity_keys, err := a.keysApp.GetUserKeys(user.ID, "bundle")
	if err != nil {
		return nil, "", "", err
	}

	encrypted_master_key, err := a.keysApp.GetUserKeys(user.ID, "master")
	if err != nil {
		return nil, "", "", err
	}

	return &domain.KeysRequest{
		IdentityKeys: domain.IdentityKeysRequest{
			EncryptedSecretKeys: domain.EncryptedKey{
				Ciphertext: encrypted_identity_keys.Ciphertext,
				Nonce:      encrypted_identity_keys.Nonce,
				Salt:       encrypted_identity_keys.Salt,
				Signature:  encrypted_identity_keys.Signature,
			},
			IdentityPublicKeys: domain.IdentityPublicKeys{
				MlKemPublicKey: user.MlKemPublicKey,
				EcdhPublicKey:  user.EcdhPublicKey,
				EdPublicKey:    user.EdPublicKey,
			},
		},
		EncryptedMasterKey: domain.MasterKeyRequest{
			EncryptedKey: domain.EncryptedKey{
				Ciphertext: encrypted_master_key.Ciphertext,
				Nonce:      encrypted_master_key.Nonce,
				Salt:       encrypted_master_key.Salt,
				Signature:  encrypted_master_key.Signature,
			},
		},
	}, challenge, user.PublicID, nil
}
