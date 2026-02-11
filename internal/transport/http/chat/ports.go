package chat

import "github.com/slipe-fun/skid-backend/internal/domain"

type ChatApp interface {
	CreateChat(token string, recipient int) (*domain.Chat, *domain.Session, error)
	GetChatByID(token string, id int) (*domain.Chat, error)
	GetChatsByUserID(token string) ([]*domain.ChatWithLastMessage, error)
	GetChatWithUsers(token string, recipient int) (*domain.Chat, error)
	GetOtherMember(chat *domain.Chat, memberID int) *domain.Member
	AddKeys(token string, chat *domain.Chat, kyberPublicKey string, ecdhPublicKey string, edPublicKey string) error
}

type MessageApp interface {
	GetChatMessages(token string, chatID int) ([]*domain.MessageWithReply, error)
	GetChatMessagesAfter(token string, chatID int, afterID int, count int) ([]*domain.MessageWithReply, error)
	GetChatMessagesBefore(token string, chatID, beforeID, count int) ([]*domain.MessageWithReply, error)
	GetChatLastReadMessage(token string, chatID int) (*domain.Message, error)
}

type UserApp interface {
	GetUserByID(id int) (*domain.User, error)
}
