package user

import "github.com/slipe-fun/skid-backend/internal/domain"

type UserRepo interface {
	GetByID(id int) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetAllUsers(limit, offset int) ([]*domain.User, error)
	Edit(user *domain.User) error
	SearchUsersByUsername(query string, limit, offset int) ([]*domain.User, error)
	UpdatePublicKeys(userID int, kyber, ecdh, ed string) error
}

type KeysRepo interface {
	Create(keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error)
}
