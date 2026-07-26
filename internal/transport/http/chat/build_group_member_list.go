package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto/validations"
)

type Member struct {
	MemberID          string                   `json:"member_id"`
	Handshake         domain.Handshake         `json:"handshake"`
	EncryptedGroupKey domain.EncryptedGroupKey `json:"encrypted_group_key"`
}

type MembersRequest struct {
	Members []Member `json:"members"`
}

type Error struct {
	Status   int
	Response fiber.Map
}

func (h *ChatHandler) BuildGroupMemberList(req MembersRequest, sessionUser *domain.User) ([]domain.GroupMember, *Error) {
	var members []domain.User
	users := make(map[string]domain.User, len(req.Members))
	var err error
	if len(req.Members) != 0 {
		var membersIDs []string
		for i := range req.Members {
			member := req.Members[i]
			membersIDs = append(membersIDs, member.MemberID)
		}

		members, err = h.userApp.GetUsersByPublicIDs(membersIDs)
		if appErr, ok := err.(*domain.AppError); ok {
			return nil, &Error{
				Status: appErr.Status,
				Response: fiber.Map{
					"error":   appErr.Code,
					"message": appErr.Msg,
				},
			}
		}

		if len(members) == 0 {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "invalid_request",
					"message": "invalid members list",
				},
			}
		}

		for i := range members {
			member := members[i]
			users[member.PublicID] = member
		}
	}

	var groupMembers []domain.GroupMember

	for i := range req.Members {
		memberObject := req.Members[i]

		member, exists := users[memberObject.MemberID]
		if !exists {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "user_not_found",
					"message": "member " + memberObject.MemberID + " not found",
				},
			}
		}

		if member.ID == sessionUser.ID {
			continue
		}

		if memberObject.Handshake.ReceiverCipherText == "" {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "no_receiver_cipher_text",
					"message": memberObject.MemberID + ": receiver cipher text is missing",
				},
			}
		}

		if memberObject.Handshake.SenderEphemeralX448 == "" {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "no_sender_cipher_text",
					"message": memberObject.MemberID + ": sender cipher text is missing",
				},
			}
		}

		if memberObject.Handshake.EncryptedSyncKey.CipherText == "" {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "no_sync_key_ciphertext",
					"message": memberObject.MemberID + ": encrypted sync key ciphertext is missing",
				},
			}
		}

		if memberObject.Handshake.EncryptedSyncKey.Nonce == "" {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "no_sync_key_nonce",
					"message": memberObject.MemberID + ": encrypted sync key nonce is missing",
				},
			}
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.ReceiverCipherText, MLKEM768CiphertextSize); err != nil {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "invalid_receiver_ciphertext_length",
					"message": memberObject.MemberID + ": invalid receiver ciphertext length",
				},
			}
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.SenderEphemeralX448, X448PublicKeySize); err != nil {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "invalid_sender_ciphertext_length",
					"message": memberObject.MemberID + ": invalid sender ciphertext length",
				},
			}
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.EncryptedSyncKey.Nonce, AESGCMNonceSize); err != nil {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "invalid_nonce_length",
					"message": memberObject.MemberID + ": invalid nonce length",
				},
			}
		}

		if err := validations.ValidateCryptoLength(memberObject.Handshake.EncryptedSyncKey.CipherText, SyncKeyCiphertextSize); err != nil {
			return nil, &Error{
				Status: fiber.ErrBadRequest.Code,
				Response: fiber.Map{
					"error":   "invalid_sync_key_length",
					"message": memberObject.MemberID + ": invalid sync key length",
				},
			}
		}

		newGroupMember := domain.GroupMember{
			MemberID:          member.ID,
			Role:              "member",
			Handshake:         memberObject.Handshake,
			EncryptedGroupKey: memberObject.EncryptedGroupKey,
		}

		groupMembers = append(groupMembers, newGroupMember)
	}

	return groupMembers, nil
}
