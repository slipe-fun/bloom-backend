package message

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

type SessionApp interface {
	GetSession(token string) (*domain.Session, error)
}

type MessageRepo interface {
	Create(message *domain.Message) (*domain.Message, error)
	UpdateMessagesSeenStatus(messages []int, seenTime time.Time) error
	GetChatLastReadMessage(chatID int) (*domain.Message, error)
	GetChatMessagesAfter(chatID, afterID, count int) ([]*domain.Message, error)
	GetChatMessagesBefore(chatID, beforeID, count int) ([]*domain.Message, error)
	GetChatMessages(id int) ([]*domain.Message, error)
	GetByID(id int) (*domain.Message, error)
}

type ChatApp interface {
	GetChatByID(user_id int, id int) (*domain.Chat, error)
	HasMember(chat *domain.Chat, memberID int) bool
}
