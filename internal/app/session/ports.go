package session

import "github.com/slipe-fun/skid-backend/internal/domain"

type SessionRepo interface {
	Create(session *domain.Session) (*domain.Session, error)
	AddKeys(id int, identity_pub, ecdh_pub, kyber_pub string) error
	GetByID(id int) (*domain.Session, error)
	GetByToken(token string) (*domain.Session, error)
	GetByUserID(id int) ([]*domain.Session, error)
	Delete(id int) error
}

type UserRepo interface {
	GetByID(id int) (*domain.User, error)
}

type JWTService interface {
	GenerateToken(userID int) (string, error)
}

type TokenService interface {
	ExtractUserID(token string) (int, error)
}
