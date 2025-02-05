package main

import (
	"finances-backend/models"
	"finances-backend/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	authService   services.AuthService
	userService   services.UserService
	walletService services.WalletService
	validate      *validator.Validate
}

func NewHandler(
	as services.AuthService,
	us services.UserService,
	ws services.WalletService,
) *Handler {
	return &Handler{
		authService:   as,
		userService:   us,
		walletService: ws,
		validate:      validator.New(),
	}
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
		c.JSON(fiber.Map{"error": "user already exists"})
		return c.SendStatus(fiber.StatusBadRequest)
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
	userID := c.Locals("user_id").(int64)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(user)
}

func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	// parse wallet's data
	wallet := models.Wallet{}
	if err := c.BodyParser(&wallet); err != nil {
		c.JSON(fiber.Map{"error": "wrong format"})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	// validate wallet
	if err := h.validate.Struct(wallet); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet.OwnerID = userID
	id, err := h.walletService.CreateWallet(&wallet)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": id})
}

func (h *Handler) GetWallets(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	wallets, err := h.walletService.GetAllWallets(userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(wallets)
}

func (h *Handler) GetWalletByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	walletID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.JSON(fiber.Map{"error": "wrong format"})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet, err := h.walletService.GetWalletByID(int64(walletID), userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(wallet)
}