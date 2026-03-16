package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type UserSessionsResponse struct {
	UserID   int               `json:"user_id"`
	Sessions []*domain.Session `json:"sessions"`
}

func (h *SessionHandler) GetUserKeyBundle(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	type query struct {
		IDs []int `query:"ids"`
	}

	var q query
	if err := c.QueryParser(&q); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	if len(q.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "ids required",
		})
	}

	for _, userID := range q.IDs {
		if session.UserID == userID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid_params",
				"message": "invalid request params",
			})
		}
	}

	sessions, err := h.sessionApp.GetSessionByUserIDs(q.IDs)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	chats, err := h.chatsRepo.GetExistingChatUsers(session.UserID, q.IDs)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	if len(chats) != len(q.IDs) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "no chat exists with one or more of the specified users",
		})
	}

	for i := range sessions {
		sessions[i].Token = ""
	}

	grouped := make(map[int][]*domain.Session)
	for _, s := range sessions {
		grouped[s.UserID] = append(grouped[s.UserID], s)
	}

	var resp []UserSessionsResponse
	for _, id := range q.IDs {
		resp = append(resp, UserSessionsResponse{
			UserID:   id,
			Sessions: grouped[id],
		})
	}

	return c.JSON(sessions)

}
