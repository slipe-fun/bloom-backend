package domain

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type Credential struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	CredentialID    []byte    `db:"credential_id" json:"credential_id"`
	PublicKey       []byte    `db:"public_key" json:"public_key"`
	AttestationType string    `db:"attestation_type" json:"attestation_type"`
	SignCount       uint32    `db:"sign_count" json:"sign_count"`
	CloneWarning    bool      `db:"clone_warning" json:"clone_warning"`
	Transport       string    `db:"transport" json:"transport"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func (c *Credential) ToWebAuthn() webauthn.Credential {
	var transports []protocol.AuthenticatorTransport
	if c.Transport != "" {
		var list []string
		_ = json.Unmarshal([]byte(c.Transport), &list)
		for _, t := range list {
			transports = append(transports, protocol.AuthenticatorTransport(t))
		}
	}
	return webauthn.Credential{
		ID:              c.CredentialID,
		PublicKey:       c.PublicKey,
		AttestationType: c.AttestationType,
		Transport:       transports,
		Authenticator: webauthn.Authenticator{
			SignCount:    c.SignCount,
			CloneWarning: c.CloneWarning,
		},
	}
}

type WebAuthnUser struct {
	User        *User
	Credentials []webauthn.Credential
}

func NewWebAuthnUser(user *User, creds []webauthn.Credential) *WebAuthnUser {
	return &WebAuthnUser{
		User:        user,
		Credentials: creds,
	}
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, uint64(u.User.ID))
	return idBytes
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.User.Username
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	if u.User.DisplayName != nil {
		return *u.User.DisplayName
	}
	return u.User.Username
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}
