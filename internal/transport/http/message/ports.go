package message

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
)

type ChatApp interface {
	GetOtherMember(chat *domain.Chat, memberID int) *domain.Member
}

type MessageApp interface {
	Send(token string, encryptionType string, message *domain.SocketMessage) (*domain.MessageWithReply, *domain.Chat, *domain.Session, error)
	GetMessageByID(token string, id int) (*domain.MessageWithReply, error)
	UpdateMessagesSeenStatus(token string, chatID int, messageIDs []int) (*[]int, *time.Time, *domain.Chat, *domain.Session, error)
}
