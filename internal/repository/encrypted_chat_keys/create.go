package keys

import (
	"fmt"
	"strings"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *EncryptedChatKeysRepo) Create(keys []*domain.EncryptedChatKeys) ([]*domain.EncryptedChatKeys, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	start := time.Now()

	valueStrings := make([]string, 0, len(keys))
	valueArgs := make([]interface{}, 0, len(keys)*9)

	for i, k := range keys {
		base := i * 9

		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				base+1, base+2, base+3,
				base+4, base+5, base+6,
				base+7, base+8, base+9,
			),
		)

		valueArgs = append(valueArgs,
			k.ChatID,
			k.SessionID,
			k.FromSessionID,
			k.EncryptedKey,
			k.EncapsulatedKey,
			k.CekWrap,
			k.CekWrapIV,
			k.Salt,
			k.Nonce,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO encrypted_chat_keys
		(chat_id, session_id, from_session_id,
		 encrypted_key, encapsulated_key,
		 cek_wrap, cek_wrap_iv,
		 salt, nonce)
		VALUES %s
		RETURNING 
			id, chat_id, session_id, from_session_id,
			encrypted_key, encapsulated_key,
			cek_wrap, cek_wrap_iv,
			salt, nonce,
			created_at
	`, strings.Join(valueStrings, ","))

	rows, err := r.db.Query(query, valueArgs...)
	duration := time.Since(start)
	metrics.ObserveDB("encrypted_chat_keys_create_many", duration, err)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var created []*domain.EncryptedChatKeys

	for rows.Next() {
		var k domain.EncryptedChatKeys
		err := rows.Scan(
			&k.ID,
			&k.ChatID,
			&k.SessionID,
			&k.FromSessionID,
			&k.EncryptedKey,
			&k.EncapsulatedKey,
			&k.CekWrap,
			&k.CekWrapIV,
			&k.Salt,
			&k.Nonce,
			&k.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		created = append(created, &k)
	}

	return created, rows.Err()
}
