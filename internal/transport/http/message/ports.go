package message

import (
	"context"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

type ChatApp interface {
	GetOtherMember(chat *domain.Chat, memberID int) *domain.User
	GetGroupMembers(ctx context.Context, groupID int) ([]domain.GroupMember, error)
}

type MessageApp interface {
	Send(ctx context.Context, user_id int, message *domain.SocketMessage) (*domain.MessageWithReply, *domain.Chat, error)
	GetMessageByID(user_id, id int) (*domain.MessageWithReply, error)
	UpdateMessagesSeenStatus(user_id, chatID int, messageIDs []int) (*[]int, *time.Time, *domain.Chat, error)
}
