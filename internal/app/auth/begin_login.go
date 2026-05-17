package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *AuthApp) BeginLogin() (string, *protocol.CredentialAssertion, error) {
	options, sessionData, err := s.webauthn.BeginDiscoverableLogin()
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, domain.Failed("failed to begin login: " + err.Error())
	}

	token := uuid.New().String()

	sessBytes, err := json.Marshal(domain.LoginSession{
		SessionData: *sessionData,
		UserID:      0,
	})
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, domain.Failed("failed to save login session")
	}

	ctx := context.Background()
	err = s.rdb.Set(ctx, "webauthn:login:"+token, sessBytes, 5*time.Minute).Err()
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", nil, domain.Failed("failed to save login session")
	}

	return token, options, nil
}
