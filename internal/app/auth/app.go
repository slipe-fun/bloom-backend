package auth

type AuthApp struct {
	sessionApp SessionApp
}

func NewAuthApp(sessionApp SessionApp) *AuthApp {
	return &AuthApp{
		sessionApp: sessionApp,
	}
}
