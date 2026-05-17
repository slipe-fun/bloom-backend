package domain

import "github.com/go-webauthn/webauthn/webauthn"

type LoginSession struct {
	SessionData webauthn.SessionData `json:"session_data"`
	UserID      int                  `json:"user_id"`
}

type RegSession struct {
	SessionData webauthn.SessionData `json:"session_data"`
	Username    string               `json:"username"`
	UserID      int                  `json:"user_id"`
}
