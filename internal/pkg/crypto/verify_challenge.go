package crypto

import (
	"encoding/json"

	"github.com/cloudflare/circl/sign/ed448"
)

type LoginChallenge struct {
	Challenge string `json:"challenge"`
	UserID    string `json:"user_id"`
}

func VerifyLoginChallengeSignature(pubKey []byte, signature []byte, challenge, userID string) (bool, error) {
	data := LoginChallenge{
		Challenge: challenge,
		UserID:    userID,
	}

	message, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	isValid := ed448.Verify(pubKey, message, signature, "")

	return isValid, nil
}
