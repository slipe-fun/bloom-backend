package session

import (
	SessionApp "github.com/slipe-fun/skid-backend/internal/app/session"
)

type SessionHandler struct {
	sessionApp *SessionApp.SessionApp
}

func NewSessionHandler(sessionApp *SessionApp.SessionApp) *SessionHandler {
	return &SessionHandler{
		sessionApp: sessionApp,
	}
}
