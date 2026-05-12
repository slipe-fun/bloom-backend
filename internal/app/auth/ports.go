package auth

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type SessionApp interface {
	CreateSession(user_id int) (string, *domain.Session, error)
}