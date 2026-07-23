package auth

import (
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/generator"
)

func (a *AuthApp) Register(req *domain.KeysRequest) (string, *domain.User, *domain.Session, error) {
	ecdhBytes, err := base64.StdEncoding.DecodeString(req.IdentityKeys.IdentityPublicKeys.EcdhPublicKey)
	if err != nil {
		return "", nil, nil, domain.InvalidData("invalid ecdh public key format")
	}

	mlKemBytes, err := base64.StdEncoding.DecodeString(req.IdentityKeys.IdentityPublicKeys.MlKemPublicKey)
	if err != nil {
		return "", nil, nil, domain.InvalidData("invalid ml-kem public key format")
	}

	user, err := a.users.GetByAuthLookupID(req.AuthLookupID)
	if user != nil {
		return "", nil, nil, domain.InvalidData("user with this auth lookup id is already exists")
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", nil, nil, domain.Failed("failed to create user")
	}

	publicID := generator.GenerateUserID(ecdhBytes, mlKemBytes)
	username, err := generator.GenerateUsername()
	if err != nil {
		return "", nil, nil, domain.Failed("failed to generate username")
	}

	createdUser, err := a.users.Create(&domain.User{
		AuthLookupID: req.AuthLookupID,
		PublicID:     publicID,
		Username:     username,
	})
	if err != nil {
		return "", nil, nil, domain.Failed("failed to create user")
	}

	rollback := func() {
		_ = a.users.Delete(createdUser.ID)
	}

	_, _, err = a.keysApp.UploadIdentityKeys(createdUser.ID, &req.IdentityKeys)
	if err != nil {
		rollback()
		return "", nil, nil, err
	}

	masterKeyPayload := &domain.EncryptedKeys{
		UserID: createdUser.ID,
		Type:   "master",
		EncryptedKey: domain.EncryptedKey{
			Ciphertext: req.EncryptedMasterKey.Ciphertext,
			Nonce:      req.EncryptedMasterKey.Nonce,
			Salt:       req.EncryptedMasterKey.Salt,
			Signature:  req.EncryptedMasterKey.Signature,
		},
	}

	_, err = a.keysApp.UploadMasterKey(createdUser.ID, masterKeyPayload)
	if err != nil {
		rollback()
		return "", nil, nil, err
	}

	finalUser, err := a.users.GetByID(createdUser.ID)
	if err != nil {
		return "", nil, nil, domain.Failed("failed to fetch registered user")
	}

	token, session, err := a.sessionApp.CreateSession(finalUser.ID)
	if err != nil {
		return "", nil, nil, err
	}

	return token, finalUser, session, nil
}
