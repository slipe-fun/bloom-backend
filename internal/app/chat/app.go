package chat

type ChatApp struct {
	chats    ChatRepo
	groups   GroupMemberRepo
	messages MessageRepo
}

func NewChatApp(chats ChatRepo, groups GroupMemberRepo, messages MessageRepo) *ChatApp {
	return &ChatApp{
		chats:    chats,
		groups:   groups,
		messages: messages,
	}
}
