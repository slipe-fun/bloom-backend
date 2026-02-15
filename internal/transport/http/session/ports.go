package session

import "github.com/slipe-fun/skid-backend/internal/domain"

type SessionApp interface {
	DeleteSession(user_id, id int) error
	GetSession(token string) (*domain.Session, error)
	GetUserSessions(user_id int) ([]*domain.Session, error)
}

type ChatsRepo interface {
	GetWithUsers(id, recipient int) (*domain.Chat, error)
}