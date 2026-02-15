package session

type SessionHandler struct {
	sessionApp SessionApp
	chatsRepo  ChatsRepo
}

func NewSessionHandler(sessionApp SessionApp, chatsRepo ChatsRepo) *SessionHandler {
	return &SessionHandler{
		sessionApp: sessionApp,
		chatsRepo:  chatsRepo,
	}
}
