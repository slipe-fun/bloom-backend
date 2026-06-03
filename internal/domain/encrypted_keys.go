package domain

type EncryptedKeyBytes struct {
	Ciphertext []byte
	Nonce      []byte
	Salt       []byte
	Signature  []byte
}

type IdentityPublicKeysBytes struct {
	MlKemPublicKey []byte
	EcdhPublicKey  []byte
	EdPublicKey    []byte
}

type EncryptedKey struct {
	Ciphertext string `db:"ciphertext" json:"ciphertext"`
	Nonce      string `db:"nonce" json:"nonce"`
	Salt       string `db:"salt" json:"salt"`
	Signature  string `db:"signature" json:"signature"`
}

type IdentityPublicKeys struct {
	MlKemPublicKey string `json:"ml_kem_public_key"`
	EcdhPublicKey  string `json:"ecdh_public_key"`
	EdPublicKey    string `json:"ed_public_key"`
}

type EncryptedKeys struct {
	EncryptedKey
	ID     int    `db:"id" json:"id"`
	Type   string `db:"type" json:"type"`
	UserID int    `db:"user_id" json:"user_id"`
}

type IdentityKeysRequest struct {
	EncryptedSecretKeys EncryptedKey       `json:"encrypted_secret_keys"`
	IdentityPublicKeys  IdentityPublicKeys `json:"public_keys"`
}

type MasterKeyRequest struct {
	EncryptedKey
}

type KeysRequest struct {
	EncryptedIdentityKeys IdentityKeysRequest `json:"identity_keys"`
	EncryptedMasterKey    MasterKeyRequest    `json:"encrypted_master_key"`
}
