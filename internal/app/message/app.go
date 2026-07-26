package message

type MessageApp struct {
	messages MessageRepo
	groups   GroupMemberRepo
	chats    ChatApp
}

func NewMessageApp(
	messages MessageRepo,
	groups GroupMemberRepo,
	chats ChatApp,
) *MessageApp {
	return &MessageApp{
		messages: messages,
		groups:   groups,
		chats:    chats,
	}
}
