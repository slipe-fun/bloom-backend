package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *AuthApp) FinishRegistration(token string, responseBytes []byte) (string, *domain.Session, *domain.User, error) {
	ctx := context.Background()
	sessKey := "webauthn:register:" + token

	sessBytes, err := s.rdb.Get(ctx, sessKey).Bytes()
	if err != nil {
		return "", nil, nil, domain.Expired("registration session expired or not found")
	}

	var regSess domain.RegSession
	if err := json.Unmarshal(sessBytes, &regSess); err != nil {
		return "", nil, nil, domain.InvalidData("invalid session data")
	}

	httpReq, err := http.NewRequest("POST", "/", bytes.NewReader(responseBytes))
	if err != nil {
		return "", nil, nil, domain.Failed("failed to parse response")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	createdUser, err := s.users.GetByID(regSess.UserID)
	if err != nil || createdUser == nil {
		logger.LogError("user not found", "auth-app")
		return "", nil, nil, domain.Failed("failed to find user in database")
	}

	webauthnUser := domain.NewWebAuthnUser(createdUser, []webauthn.Credential{})

	credential, err := s.webauthn.FinishRegistration(webauthnUser, regSess.SessionData, httpReq)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.InvalidData("credential verification failed: " + err.Error())
	}

	var transportStrings []string
	for _, t := range credential.Transport {
		transportStrings = append(transportStrings, string(t))
	}
	transportBytes, _ := json.Marshal(transportStrings)

	dbCred := &domain.Credential{
		UserID:          createdUser.ID,
		CredentialID:    credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		SignCount:       credential.Authenticator.SignCount,
		CloneWarning:    credential.Authenticator.CloneWarning,
		Transport:       string(transportBytes),
	}

	_, err = s.credentials.Create(dbCred)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.Failed("failed to store credential")
	}

	s.rdb.Del(ctx, sessKey)

	sessionToken, session, err := s.sessionApp.CreateSession(createdUser.ID)
	if err != nil {
		return "", nil, nil, err
	}

	return sessionToken, session, createdUser, nil
}
