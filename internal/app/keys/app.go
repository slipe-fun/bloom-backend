package keys

type KeysApp struct {
	keys  KeysRepo
	users UserRepo
}

func NewKeysApp(keys KeysRepo, users UserRepo) *KeysApp {
	return &KeysApp{
		keys:  keys,
		users: users,
	}
}
