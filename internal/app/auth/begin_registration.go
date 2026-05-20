package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *AuthApp) BeginRegistration() (string, string, *protocol.CredentialCreation, error) {
	adjectives := []string{
		"silent", "stealth", "shadow", "cyber", "cosmic", "quantum", "neon", "rapid", "hyper", "crystal",
		"phantom", "solar", "lunar", "stellar", "polar", "golden", "silver", "crimson", "azure", "emerald",
	}
	nouns := []string{
		"storm", "pulse", "frost", "matrix", "wave", "flux", "vertex", "nexus", "quasar", "proton",
		"laser", "radar", "sonar", "orbit", "beacon", "vector", "zenith", "vortex", "helix", "apex",
	}

	var username string
	var exists = true

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 15; i++ {
		adj := adjectives[r.Intn(len(adjectives))]
		noun := nouns[r.Intn(len(nouns))]
		num := r.Intn(9000) + 1000 // 4 digits
		candidate := fmt.Sprintf("%s-%s-%d", adj, noun, num)

		existingUser, err := s.users.GetByUsername(candidate)
		if err != nil || existingUser == nil {
			username = candidate
			exists = false
			break
		}
	}

	if exists {
		return "", "", nil, domain.Failed("failed to generate unique username")
	}

	newUser := &domain.User{
		Username: username,
		Date:     time.Now(),
	}
	createdUser, err := s.users.Create(newUser)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", "", nil, domain.Failed("failed to create user in database")
	}

	webauthnUser := domain.NewWebAuthnUser(createdUser, []webauthn.Credential{})

	options, sessionData, err := s.webauthn.BeginRegistration(
		webauthnUser,
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
	)
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", "", nil, domain.Failed("failed to begin registration")
	}

	token := uuid.New().String()

	sessBytes, err := json.Marshal(domain.RegSession{
		SessionData: *sessionData,
		Username:    username,
		UserID:      createdUser.ID,
	})
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", "", nil, domain.Failed("failed to save registration session")
	}

	ctx := context.Background()
	err = s.rdb.Set(ctx, "webauthn:register:"+token, sessBytes, 5*time.Minute).Err()
	if err != nil {
		logger.LogError(err.Error(), "auth-app")
		return "", "", nil, domain.Failed("failed to save registration session")
	}

	return token, username, options, nil
}
