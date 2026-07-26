package chat

import (
	"context"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

type ChatApp interface {
	CreatePrivateChat(user_id, recipient int, handshake domain.Handshake) (*domain.Chat, error)
	CreateGroupChat(ctx context.Context, title string, members []domain.GroupMember, invitedByID int) (*domain.Chat, error)
	GetChatByID(ctx context.Context, user_id int, id int) (*domain.Chat, error)
	GetGroupMember(ctx context.Context, groupID, memberID int) (*domain.GroupMember, error)
	GetChatsByUserID(user_id int) ([]*domain.ChatWithLastMessage, error)
	GetChatWithUsers(user_id, recipient int) (*domain.Chat, error)
	GetMemberGroups(ctx context.Context, memberID int) ([]domain.GroupMember, error)
	GetOtherMember(chat *domain.Chat, memberID int) *domain.User
}

type MessageApp interface {
	GetChatMessagesAfter(ctx context.Context, user_id, chatID, afterID, count int) ([]*domain.MessageWithReply, error)
	GetChatMessagesBefore(ctx context.Context, user_id, chatID, beforeID, count int) ([]*domain.MessageWithReply, error)
	GetChatLastReadMessage(ctx context.Context, user_id, chatID int) (*domain.Message, error)
}

type UserApp interface {
	GetUserByID(id int) (*domain.User, error)
	GetUsersByIDs(ids []int) ([]domain.User, error)
	GetUserByPublicID(id string) (*domain.User, error)
	GetUsersByPublicIDs(ids []string) ([]domain.User, error)
}
