package groupmember

import (
	"context"
	"encoding/json"
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *GroupMemberRepo) Create(ctx context.Context, chatID int, invitedByID int, member *domain.GroupMember) error {
	start := time.Now()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		metrics.ObserveDB("group_member_create_begin", time.Since(start), err)
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	handshakeJSON, err := json.Marshal(member.Handshake)
	if err != nil {
		return err
	}

	memberQuery := `
		INSERT INTO group_members (chat_id, member_id, invited_by_id, handshake, role)
		VALUES ($1, $2, $3, $4, $5)
	`
	role := "member"
	if member.Role != "" {
		role = member.Role
	}

	_, err = tx.ExecContext(ctx, memberQuery, chatID, member.MemberID, invitedByID, handshakeJSON, role)
	if err != nil {
		return err
	}

	keyQuery := `
		INSERT INTO group_encrypted_keys (chat_id, member_id, ciphertext, nonce, salt)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, keyQuery,
		chatID,
		member.MemberID,
		member.EncryptedGroupKey.Ciphertext,
		member.EncryptedGroupKey.Nonce,
		member.EncryptedGroupKey.Salt,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()

	duration := time.Since(start)
	metrics.ObserveDB("group_member_create", duration, err)

	return err
}

func (r *GroupMemberRepo) CreateMany(ctx context.Context, chatID int, invitedByID int, members []domain.GroupMember) error {
	if len(members) == 0 {
		return nil
	}

	start := time.Now()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		metrics.ObserveDB("group_member_create_many_begin", time.Since(start), err)
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	memberQuery := `
		INSERT INTO group_members (chat_id, member_id, invited_by_id, handshake, role)
		VALUES ($1, $2, $3, $4, $5)
	`
	keyQuery := `
		INSERT INTO group_encrypted_keys (chat_id, member_id, ciphertext, nonce, salt)
		VALUES ($1, $2, $3, $4, $5)
	`

	stmtMember, err := tx.PrepareContext(ctx, memberQuery)
	if err != nil {
		return err
	}
	defer stmtMember.Close()

	stmtKey, err := tx.PrepareContext(ctx, keyQuery)
	if err != nil {
		return err
	}
	defer stmtKey.Close()

	for i := range members {
		member := &members[i]

		handshakeJSON, errMarshal := json.Marshal(member.Handshake)
		if errMarshal != nil {
			err = errMarshal
			return err
		}

		role := "member"
		if member.Role != "" {
			role = member.Role
		}

		_, err = stmtMember.ExecContext(ctx, chatID, member.MemberID, invitedByID, handshakeJSON, role)
		if err != nil {
			return err
		}

		_, err = stmtKey.ExecContext(ctx,
			chatID,
			member.MemberID,
			member.EncryptedGroupKey.Ciphertext,
			member.EncryptedGroupKey.Nonce,
			member.EncryptedGroupKey.Salt,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	duration := time.Since(start)
	metrics.ObserveDB("group_member_create_many", duration, err)

	return err
}
