package encryptedchatkeys

type EncryptedChatKeysApp struct {
	keys    EncryptedChatKeysRepo
	chats   ChatRepo
	session SessionApp
}

func NewEncryptedChatKeysApp(keys EncryptedChatKeysRepo, chats ChatRepo, session SessionApp) *EncryptedChatKeysApp {
	return &EncryptedChatKeysApp{
		keys:    keys,
		chats:   chats,
		session: session,
	}
}
