package handlers

import (
	"finances-backend/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Register(c *fiber.Ctx) error {
	user := models.User{}
	if err := h.parseBody(c, &user); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := h.authService.CreateUser(&user)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	c.JSON(fiber.Map{"id": id})
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	logReq := models.User{}
	if err := h.parseBody(c, &logReq); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Authenticate user
	token, err := h.authService.AuthenticateUser(&logReq)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrAuthFailed.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.JSON(fiber.Map{"token": token})
}

func (h *Handler) AuthorizeMiddleware(c *fiber.Ctx) error {
	payload, err := h.authService.AuthorizeUser(c.Locals("user"))
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrAuthFailed.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("user_id", payload)
	return c.Next()
}
