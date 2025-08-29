package app

import (
	"errors"

	"github.com/slipe-fun/skid-backend/internal/domain"
	ChatRepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
	"github.com/slipe-fun/skid-backend/internal/service"
)

type ChatApp struct {
	chats    *ChatRepo.ChatRepo
	tokenSvc *service.TokenService
}

func NewChatApp(chats *ChatRepo.ChatRepo, tokenSvc *service.TokenService) *ChatApp {
	return &ChatApp{
		chats:    chats,
		tokenSvc: tokenSvc,
	}
}

func (c *ChatApp) HasMember(chat *domain.Chat, memberID int) bool {
	for _, m := range chat.Members {
		if m.ID == memberID {
			return true
		}
	}
	return false
}

func (c *ChatApp) GetChatById(tokenStr string, id int) (*domain.Chat, error) {
	userID, err := c.tokenSvc.ExtractUserID(tokenStr)

	if err != nil {
		return nil, err
	}

	chat, err := c.chats.GetById(id)

	if err != nil {
		return nil, err
	}

	if !c.HasMember(chat, userID) {
		return nil, errors.New("user is not in chat")
	}

	return chat, nil
}

func (c *ChatApp) GetChatsByUserId(tokenStr string) ([]*domain.Chat, error) {
	userID, err := c.tokenSvc.ExtractUserID(tokenStr)
	if err != nil {
		return nil, err
	}

	chats, err := c.chats.GetByUserId(userID)

	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (c *ChatApp) GetChatWithUsers(tokenStr string, recipient int) (*domain.Chat, error) {
	userID, err := c.tokenSvc.ExtractUserID(tokenStr)
	if err != nil {
		return nil, err
	}

	chat, err := c.chats.GetWithUsers(userID, recipient)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *ChatApp) CreateChat(tokenStr string, recipient int) (*domain.Chat, error) {
	userID, err := c.tokenSvc.ExtractUserID(tokenStr)
	if err != nil {
		return nil, err
	}

	chat, err := c.chats.Create(&domain.Chat{
		Members: []domain.Member{
			{
				ID:             userID,
				KyberPublicKey: "",
				EcdhPublicKey:  "",
				EdPublicKey:    "",
			},
			{
				ID:             recipient,
				KyberPublicKey: "",
				EcdhPublicKey:  "",
				EdPublicKey:    "",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return chat, nil
}
