package encryptedchatkeys

import "github.com/slipe-fun/skid-backend/internal/transport/ws/types"

type EncryptedChatKeysHandler struct {
	keys  EncryptedChatKeysApp
	wsHub *types.Hub
}

func NewEncryptedChatKeysHandlerApp(keys EncryptedChatKeysApp, wsHub *types.Hub) *EncryptedChatKeysHandler {
	return &EncryptedChatKeysHandler{
		keys:  keys,
		wsHub: wsHub,
	}
}
