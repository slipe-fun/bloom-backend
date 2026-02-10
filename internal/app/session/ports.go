package session

import "github.com/slipe-fun/skid-backend/internal/domain"

type SessionRepo interface {
	Create(session *domain.Session) (*domain.Session, error)
	GetById(id int) (*domain.Session, error)
	GetByToken(token string) (*domain.Session, error)
	GetByUserId(id int) ([]*domain.Session, error)
	Delete(id int) error
}

type UserRepo interface {
	GetById(id int) (*domain.User, error)
}

type JWTService interface {
	GenerateToken(userID int) (string, error)
}

type TokenService interface {
	ExtractUserID(tokenStr string) (int, error)
}
