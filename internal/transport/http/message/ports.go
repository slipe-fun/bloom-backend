package message

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

type ChatApp interface {
	GetOtherMember(chat *domain.Chat, memberID int) *domain.Member
}

type MessageApp interface {
	Send(user_id int, encryptionType domain.EncryptionType, message *domain.SocketMessage) (*domain.MessageWithReply, *domain.Chat, error)
	GetMessageByID(user_id, id int) (*domain.MessageWithReply, error)
	UpdateMessagesSeenStatus(user_id, chatID int, messageIDs []int) (*[]int, *time.Time, *domain.Chat, error)
}
