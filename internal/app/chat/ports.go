package chat

import (
	"context"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

type ChatRepo interface {
	Create(chat *domain.RawChat) (*domain.Chat, error)
	UpdateChat(chat *domain.Chat) error
	GetByID(id int) (*domain.Chat, error)
	GetWithUsers(id, recipient int) (*domain.Chat, error)
	GetByUserID(userID int) ([]*domain.ChatWithLastMessage, error)
}

type GroupMemberRepo interface {
	CreateMany(ctx context.Context, chatID int, invitedByID int, members []domain.GroupMember) error
	GetByGroupID(ctx context.Context, chatID int) ([]domain.GroupMember, error)
	GetByMemberAndChatID(ctx context.Context, memberID int, chatID int) (*domain.GroupMember, error)
}

type MessageRepo interface {
	GetChatLastReadMessage(chatID int) (*domain.Message, error)
}
