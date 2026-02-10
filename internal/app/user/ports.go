package user

import "github.com/slipe-fun/skid-backend/internal/domain"

type SessionApp interface {
	GetSession(token string) (*domain.Session, error)
}

type UserRepo interface {
	GetById(id int) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Edit(user *domain.User) error
	SearchUsersByUsername(query string, limit, offset int) ([]*domain.User, error)
}
