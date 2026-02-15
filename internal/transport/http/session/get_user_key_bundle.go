package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (h *SessionHandler) GetUserKeyBundle(c *fiber.Ctx) error {
	sessionVal := c.Locals("session")
	session, ok := sessionVal.(*domain.Session)
	if !ok {
		return fiber.ErrUnauthorized
	}

	user_id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	if session.UserID == user_id {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_params",
			"message": "invalid request params",
		})
	}

	_, err = h.chatsRepo.GetWithUsers(session.UserID, user_id)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	sessions, err := h.sessionApp.GetUserSessions(user_id)
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{
			"error":   appErr.Code,
			"message": appErr.Msg,
		})
	}

	for i := range sessions {
		sessions[i].Token = ""
	}

	return c.JSON(sessions)

}
