package groupmember

import (
	"context"
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *GroupMemberRepo) GetByGroupID(ctx context.Context, chatID int) ([]domain.GroupMember, error) {
	start := time.Now()

	query := `
		SELECT
			m.chat_id, m.member_id, m.invited_by_id, m.role, m.joined_at, m.handshake,
			k.ciphertext, k.nonce, k.salt
		FROM group_members m
		LEFT JOIN group_encrypted_keys k
			ON m.chat_id = k.chat_id AND m.member_id = k.member_id
		WHERE m.chat_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, chatID)
	if err != nil {
		metrics.ObserveDB("group_member_get_by_group_id", time.Since(start), err)
		return nil, err
	}
	defer rows.Close()

	var members []domain.GroupMember

	for rows.Next() {
		var m domain.GroupMember
		var handshakeBytes []byte

		err := rows.Scan(
			&m.ChatID,
			&m.MemberID,
			&m.InvitedByID,
			&m.Role,
			&m.JoinedAt,
			&handshakeBytes,
			&m.EncryptedGroupKey.Ciphertext,
			&m.EncryptedGroupKey.Nonce,
			&m.EncryptedGroupKey.Salt,
		)
		if err != nil {
			return nil, err
		}

		if len(handshakeBytes) > 0 {
			if err := json.Unmarshal(handshakeBytes, &m.Handshake); err != nil {
				return nil, err
			}
		}

		members = append(members, m)
	}

	err = rows.Err()
	metrics.ObserveDB("group_member_get_by_group_id", time.Since(start), err)

	if err != nil {
		return nil, err
	}

	return members, nil
}
