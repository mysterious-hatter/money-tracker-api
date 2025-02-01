package main

import (
	"finances-backend/models"
	"finances-backend/services"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	userService services.UserService
	validate *validator.Validate
}

func NewHandler(userService services.UserService) *Handler {
	return &Handler{userService: userService, validate: validator.New()}
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

	id, err := h.userService.CreateUser(&user)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(fiber.Map{"id": id})
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	logReq := LoginRequest{}
	if err := c.BodyParser(&logReq); err != nil {
		c.JSON(fiber.Map{"error": "wrong format"})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.validate.Struct(logReq); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Throws Unauthorized error
	if logReq.Name != "admin" || logReq.Password != "admin" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	expiration, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"userid":  0,
		"exp":   time.Now().Add(time.Hour * time.Duration(expiration)).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(fiber.Map{"error": err})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	c.JSON(*users)
	return c.SendStatus(fiber.StatusOK)
}

type LoginRequest struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=50"`
}