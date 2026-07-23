package auth

import "github.com/slipe-fun/skid-backend/internal/domain"

type AuthApp interface {
	Register(req *domain.KeysRequest) (string, *domain.User, *domain.Session, error)
	LoginBegin(authLookupID string) (*domain.KeysRequest, string, string, error)
	LoginFinish(user_id, signature string) (string, *domain.User, *domain.Session, error)
}
