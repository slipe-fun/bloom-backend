package auth

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type AuthApp interface {
	BeginRegistration() (string, string, *protocol.CredentialCreation, error)
	FinishRegistration(token string, responseBytes []byte) (string, *domain.Session, *domain.User, error)
	BeginLogin() (string, *protocol.CredentialAssertion, error)
	FinishLogin(token string, responseBytes []byte) (string, *domain.Session, *domain.User, error)
}