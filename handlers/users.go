package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotGetProfile.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(user)
}
