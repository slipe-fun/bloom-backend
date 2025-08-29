package service

import "github.com/slipe-fun/skid-backend/internal/domain"

func IsUserInChat(chats []*domain.Chat, chatId int) bool {
	allowed := false
	for _, chat := range chats {
		if chat.ID == chatId {
			allowed = true
			break
		}
	}

	return allowed
}
