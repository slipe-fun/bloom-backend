package session

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (s *SessionApp) AddKeys(id, user_id int, identity_pub, ecdh_pub, kyber_pub string) error {
	session, err := s.session.GetByID(id)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return domain.NotFound("session not found")
	}

	if session.UserID != user_id {
		// fake not found if session isnt belongs to user_id
		return domain.NotFound("session not found")
	}

	if session.IdentityPublicKey != nil || session.EcdhPublicKey != nil || session.KyberPublicKey != nil {
		return domain.Failed("keys already set for this session")
	}

	if len(identity_pub) != 32 {
		return domain.InvalidData("invalid identity public key length")
	}

	if len(ecdh_pub) != 32 {
		return domain.InvalidData("invalid ecdh public key length")
	}

	if len(kyber_pub) != 1184 {
		return domain.InvalidData("invalid kyber public key length")
	}

	err = s.session.AddKeys(id, identity_pub, ecdh_pub, kyber_pub)
	if err != nil {
		logger.LogError(err.Error(), "session-app")
		return domain.Failed("failed to add session keys")
	}

	return nil
}
