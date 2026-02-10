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
	GetChatMessagesAfter(chatId, afterId, count int) ([]*domain.Message, error)
	GetChatMessagesBefore(chatId, afterId, count int) ([]*domain.Message, error)
	GetChatMessages(id int) ([]*domain.Message, error)
	GetById(id int) (*domain.Message, error)
}

type ChatApp interface {
	GetChatById(tokenStr string, id int) (*domain.Chat, error)
	HasMember(chat *domain.Chat, memberID int) bool
}
