package domain

type EncryptionType string

const (
	ServerEncryption EncryptionType = "server"
	ClientEncryption EncryptionType = "client"
)

func (e EncryptionType) IsValid() bool {
	switch e {
	case ServerEncryption, ClientEncryption:
		return true
	default:
		return false
	}
}
