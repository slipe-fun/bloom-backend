package chat

import "github.com/slipe-fun/skid-backend/internal/domain"

type ChatApp interface {
	CreateChat(user_id, recipient int) (*domain.Chat, error)
	GetChatByID(user_id int, id int) (*domain.Chat, error)
	GetChatsByUserID(user_id int) ([]*domain.ChatWithLastMessage, error)
	GetChatWithUsers(user_id, recipient int) (*domain.Chat, error)
	GetOtherMember(chat *domain.Chat, memberID int) *domain.Member
	AddKeys(user_id int, chat *domain.Chat, kyberPublicKey string, ecdhPublicKey string, edPublicKey string) error
}

type MessageApp interface {
	GetChatMessages(user_id, chatID int) ([]*domain.MessageWithReply, error)
	GetChatMessagesAfter(user_id, chatID, afterID, count int) ([]*domain.MessageWithReply, error)
	GetChatMessagesBefore(user_id, chatID, beforeID, count int) ([]*domain.MessageWithReply, error)
	GetChatLastReadMessage(user_id, chatID int) (*domain.Message, error)
}

type UserApp interface {
	GetUserByID(id int) (*domain.User, error)
}
