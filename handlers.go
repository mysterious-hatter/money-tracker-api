package main

import (
	"finances-backend/services"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	userService services.UserService
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	return c.SendString("Users")
}