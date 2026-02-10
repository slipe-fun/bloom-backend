package chat

type ChatApp struct {
	sessionApp SessionApp
	chats      ChatRepo
}

func NewChatApp(sessionApp SessionApp, chats ChatRepo) *ChatApp {
	return &ChatApp{
		sessionApp: sessionApp,
		chats:      chats,
	}
}
