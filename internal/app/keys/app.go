package keys

type KeysApp struct {
	sessionApp SessionApp
	keys       KeysRepo
}

func NewKeysApp(sessionApp SessionApp,
	keys KeysRepo,
) *KeysApp {
	return &KeysApp{
		sessionApp: sessionApp,
		keys:       keys,
	}
}
