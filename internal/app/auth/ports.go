package auth

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type SessionApp interface {
	CreateSession(user_id int) (string, *domain.Session, error)
}

type UserRepo interface {
	GetByID(id int) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
}
