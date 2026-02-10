package message

type MessageApp struct {
	sessionApp SessionApp
	messages   MessageRepo
	chats      ChatApp
}

func NewMessageApp(sessionApp SessionApp,
	messages MessageRepo,
	chats ChatApp,
) *MessageApp {
	return &MessageApp{
		sessionApp: sessionApp,
		messages:   messages,
		chats:      chats,
	}
}
