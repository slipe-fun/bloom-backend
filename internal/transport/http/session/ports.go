package session

import "github.com/slipe-fun/skid-backend/internal/domain"

type SessionApp interface {
	AddKeys(id, user_id int, identity_pub, ecdh_pub, kyber_pub string) error
	DeleteSession(user_id, id int) error
	GetSession(token string) (*domain.Session, error)
	GetUserSessions(user_id int) ([]*domain.Session, error)
	GetSessionByIDs(ids []int) ([]*domain.Session, error)
	GetSessionByUserIDs(ids []int) ([]*domain.Session, error)
}

type ChatsRepo interface {
	GetWithUsers(id, recipient int) (*domain.Chat, error)
	GetExistingChatUsers(userID int, ids []int) ([]int, error)
}
