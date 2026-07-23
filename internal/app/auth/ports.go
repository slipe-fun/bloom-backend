package auth

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type SessionApp interface {
	CreateSession(user_id int) (string, *domain.Session, error)
}

type KeysApp interface {
	UploadIdentityKeys(user_id int, req *domain.IdentityKeysRequest) (*domain.IdentityPublicKeysBytes, *domain.EncryptedKeyBytes, error)
	UploadMasterKey(user_id int, key *domain.EncryptedKeys) (*domain.EncryptedKeyBytes, error)
	GetUserKeys(user_id int, keys_type string) (*domain.EncryptedKeys, error)
}

type UserRepo interface {
	GetByID(id int) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByPublicID(id string) (*domain.User, error)
	GetByAuthLookupID(authLookupID string) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
	Delete(id int) error
}
