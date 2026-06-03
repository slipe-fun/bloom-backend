package auth

import "github.com/slipe-fun/skid-backend/internal/domain"

type AuthApp interface {
	Register(req *domain.KeysRequest) (string, *domain.User, *domain.Session, error)
	LoginBegin(user_id string) (*domain.KeysRequest, string, error)
}
