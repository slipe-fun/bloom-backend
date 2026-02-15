package auth

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"golang.org/x/oauth2"
)

type SessionApp interface {
	CreateSession(user_id int) (string, *domain.Session, error)
}

type UserRepo interface {
	GetByEmail(email string) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
}

type VerificationApp interface {
	CreateAndSendCode(email string) error
}

type VerificationRepo interface {
	GetLastCode(email string) (*domain.VerificationCode, error)
	GetByEmailAndCode(email, code string) (*domain.VerificationCode, error)
	DeleteByEmailAndCode(email, code string) error
}

type GoogleAuthService interface {
	ExchangeCode(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (map[string]interface{}, error)
}
