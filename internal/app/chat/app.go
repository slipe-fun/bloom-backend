package chat

type ChatApp struct {
	chats    ChatRepo
	messages MessageRepo
}

func NewChatApp(chats ChatRepo, messages MessageRepo) *ChatApp {
	return &ChatApp{
		chats:    chats,
		messages: messages,
	}
}
