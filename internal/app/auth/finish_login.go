package auth

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *AuthApp) FinishLogin(token string, responseBytes []byte) (string, *domain.Session, *domain.User, error) {
	ctx := context.Background()
	sessKey := "webauthn:login:" + token

	sessBytes, err := s.rdb.Get(ctx, sessKey).Bytes()
	if err != nil {
		return "", nil, nil, domain.Expired("login session expired or not found")
	}

	var loginSess domain.LoginSession
	if err := json.Unmarshal(sessBytes, &loginSess); err != nil {
		return "", nil, nil, domain.InvalidData("invalid session data")
	}

	httpReq, err := http.NewRequest("POST", "/", bytes.NewReader(responseBytes))
	if err != nil {
		return "", nil, nil, domain.Failed("failed to parse response")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	parsedAssertionReq, err := http.NewRequest("POST", "/", bytes.NewReader(responseBytes))
	if err != nil {
		return "", nil, nil, domain.Failed("failed to parse response")
	}
	parsedAssertionReq.Header.Set("Content-Type", "application/json")

	parsedAssertion, err := protocol.ParseCredentialRequestResponse(parsedAssertionReq)
	if err != nil {
		return "", nil, nil, domain.InvalidData("failed to parse request response: " + err.Error())
	}

	userHandle := parsedAssertion.Response.UserHandle
	userID := 0
	if len(userHandle) >= 8 {
		userID = int(binary.BigEndian.Uint64(userHandle))
	} else if len(userHandle) > 0 {
		temp := make([]byte, 8)
		copy(temp[8-len(userHandle):], userHandle)
		userID = int(binary.BigEndian.Uint64(temp))
	}
	user, err := s.users.GetByID(userID)
	if err != nil || user == nil {
		return "", nil, nil, domain.NotFound(fmt.Sprintf("user not found (decoded userID=%d, handleLen=%d)", userID, len(userHandle)))
	}

	dbCreds, err := s.credentials.GetByUserID(user.ID)
	if err != nil {
		return "", nil, nil, domain.Failed("failed to retrieve credentials")
	}
	var webauthnCreds []webauthn.Credential
	for _, c := range dbCreds {
		webauthnCreds = append(webauthnCreds, c.ToWebAuthn())
	}

	webauthnUser := domain.NewWebAuthnUser(user, webauthnCreds)

	loginSess.SessionData.UserID = webauthnUser.WebAuthnID()

	credential, err := s.webauthn.FinishLogin(webauthnUser, loginSess.SessionData, httpReq)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.InvalidData("credential verification failed: " + err.Error())
	}

	err = s.credentials.UpdateSignCount(credential.ID, credential.Authenticator.SignCount, credential.Authenticator.CloneWarning)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, nil, domain.Failed("failed to update credential sign count")
	}

	s.rdb.Del(ctx, sessKey)

	sessionToken, session, err := s.sessionApp.CreateSession(user.ID)
	if err != nil {
		return "", nil, nil, err
	}

	return sessionToken, session, user, nil
}
