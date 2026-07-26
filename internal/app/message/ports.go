package message

import (
	"context"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

type MessageRepo interface {
	Create(message *domain.Message) (*domain.Message, error)
	UpdateMessagesSeenStatus(messages []int, seenTime time.Time) error
	GetChatLastReadMessage(chatID int) (*domain.Message, error)
	GetChatMessagesAfter(chatID, afterID, count int) ([]*domain.Message, error)
	GetChatMessagesBefore(chatID, beforeID, count int) ([]*domain.Message, error)
	GetByID(id int) (*domain.Message, error)
}

type GroupMemberRepo interface {
	GetByMemberAndChatID(ctx context.Context, memberID int, chatID int) (*domain.GroupMember, error)
}

type ChatApp interface {
	GetChatByID(ctx context.Context, user_id int, id int) (*domain.Chat, error)
	HasMember(ctx context.Context, chat *domain.Chat, memberID int) bool
}
