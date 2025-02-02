package main

import (
	"finances-backend/models"
	"finances-backend/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService services.AuthService
	userService services.UserService
	validate *validator.Validate
}

func NewHandler(as services.AuthService, us services.UserService) *Handler {
	return &Handler{authService: as, userService: us, validate: validator.New()}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		c.JSON(fiber.Map{"error": "wrong format"})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(user); err != nil {
		c.JSON(fiber.Map{"error": "wrong format"})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := h.authService.CreateUser(&user)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(fiber.Map{"id": id})
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	logReq := models.User{}
	if err := c.BodyParser(&logReq); err != nil {
		c.JSON(fiber.Map{"error": "wrong format"})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(logReq); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	
	// Authenticate user
	token, err := h.authService.AuthenticateUser(&logReq)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.JSON(fiber.Map{"token": token})
}

func (h *Handler) AuthorizeMiddleware(c *fiber.Ctx) error {
	payload, err := h.authService.AuthorizeUser(c.Locals("user"))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	
	c.Locals("user_id", payload)
	return c.Next()
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	id := c.Locals("user_id").(int64)
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(user)
}

