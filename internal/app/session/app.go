package session

type SessionApp struct {
	session  SessionRepo
	users    UserRepo
	jwtSvc   JWTService
	tokenSvc TokenService
}

func NewSessionApp(session SessionRepo,
	users UserRepo,
	jwtSvc JWTService,
	tokenSvc TokenService) *SessionApp {
	return &SessionApp{
		session:  session,
		users:    users,
		jwtSvc:   jwtSvc,
		tokenSvc: tokenSvc,
	}
}
