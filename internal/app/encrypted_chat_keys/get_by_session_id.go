package encryptedchatkeys

import "github.com/slipe-fun/skid-backend/internal/domain"

func (k *EncryptedChatKeysApp) GetBySessionID(session_id int) ([]*domain.EnrichedChatKey, error) {
	keys, err := k.keys.GetBySessionID(session_id)
	if err != nil {
		return nil, domain.Failed("failed to get keys")
	}

	var sessionIDs []int

	for key := range keys {
		sessionIDs = append(sessionIDs, keys[key].FromSessionID)
	}

	sessions, err := k.session.GetSessionByIDs(sessionIDs)
	if err != nil {
		return nil, domain.Failed("failed to get sessions")
	}

	sessionsMap := make(map[int]*domain.Session)
	for _, s := range sessions {
		sessionsMap[s.ID] = s
	}

	var enrichedKeys []*domain.EnrichedChatKey

	for _, key := range keys {
		senderSession, exists := sessionsMap[key.FromSessionID]

		var pubKeys map[string]string
		if exists {
			pubKeys = map[string]string{
				"kyber_pub":    *senderSession.KyberPublicKey,
				"ecdh_pub":     *senderSession.EcdhPublicKey,
				"identity_pub": *senderSession.IdentityPublicKey,
			}
		} else {
			pubKeys = map[string]string{}
		}

		enrichedKeys = append(enrichedKeys, &domain.EnrichedChatKey{
			EncryptedChatKeys: key,
			SenderPublicKeys:  pubKeys,
		})
	}

	return enrichedKeys, nil
}
