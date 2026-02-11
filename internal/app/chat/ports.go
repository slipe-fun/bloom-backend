package chat

import "github.com/slipe-fun/skid-backend/internal/domain"

type SessionApp interface {
	GetSession(token string) (*domain.Session, error)
}

type ChatRepo interface {
	Create(chat *domain.Chat) (*domain.Chat, error)
	UpdateChat(chat *domain.Chat) error
	GetByID(id int) (*domain.Chat, error)
	GetWithUsers(id, recipient int) (*domain.Chat, error)
	GetByUserID(userID int) ([]*domain.ChatWithLastMessage, error)
}
