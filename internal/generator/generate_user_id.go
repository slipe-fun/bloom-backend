package generator

import (
	"github.com/mr-tron/base58/base58"
	"golang.org/x/crypto/blake2b"
)

func GenerateUserID(ecdhPK []byte, mlkemPK []byte) string {
	combined := make([]byte, len(ecdhPK)+len(mlkemPK))
	copy(combined, ecdhPK)
	copy(combined[len(ecdhPK):], mlkemPK)

	hash := blake2b.Sum256(combined)

	truncated := hash[:10]

	id := base58.Encode(truncated)

	return id
}
