package auth

import "github.com/slipe-fun/skid-backend/internal/domain"

type AuthApp interface {
	Register(req *domain.KeysRequest) (string, *domain.User, *domain.Session, error)
}
