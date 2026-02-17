package encryptedchatkeys

type EncryptedChatKeysHandler struct {
	keys EncryptedChatKeysApp
}

func NewEncryptedChatKeysHandlerApp(keys EncryptedChatKeysApp) *EncryptedChatKeysHandler {
	return &EncryptedChatKeysHandler{
		keys: keys,
	}
}
