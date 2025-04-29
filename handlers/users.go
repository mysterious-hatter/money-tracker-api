package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return h.sendError(c, ErrCannotGetProfile, err)
	}

	return h.send(c, fiber.StatusOK, user)
}
