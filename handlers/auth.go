package handlers

import (
	"finances-backend/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Register(c *fiber.Ctx) error {
	user := models.User{}
	if err := h.parseBody(c, &user); err != nil {
		return h.sendError(c, ErrWrongFormat, err.Error())
	}

	id, err := h.authService.CreateUser(&user)
	if err != nil {
		return h.sendError(c, ErrCannotCreateUser, err.Error())
	}

	return h.send(c, fiber.StatusCreated, fiber.Map{"id": id})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	logReq := models.User{}
	if err := h.parseBody(c, &logReq); err != nil {
		return h.sendError(c, ErrWrongFormat, err.Error())
	}

	// Authenticate user
	token, err := h.authService.AuthenticateUser(&logReq)
	if err != nil {
		return h.sendError(c, ErrAuthFailed, err.Error())
	}

	return h.send(c, fiber.StatusOK, fiber.Map{"token": token})
}

func (h *Handler) AuthorizeMiddleware(c *fiber.Ctx) error {
	payload, err := h.authService.AuthorizeUser(c.Locals("user"))
	if err != nil {
		return h.sendError(c, ErrAuthFailed, err.Error())
	}

	c.Locals("user_id", payload)
	return c.Next()
}
