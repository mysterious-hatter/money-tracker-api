package main

import (
	"errors"
	"finances-backend/models"
	"finances-backend/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrWrongFormat error = errors.New("wrong format")
)

type Handler struct {
	authService      services.AuthService
	userService      services.UserService
	walletService    services.WalletService
	categoryService  services.CategoryService
	operationSerivce services.OperationService
	validate         *validator.Validate
}

func NewHandler(
	as services.AuthService,
	us services.UserService,
	ws services.WalletService,
	cs services.CategoryService,
	ops services.OperationService,
) *Handler {
	return &Handler{
		authService:      as,
		userService:      us,
		walletService:    ws,
		categoryService:  cs,
		operationSerivce: ops,
		validate:         validator.New(),
	}
}

// Auth
func (h *Handler) Register(c *fiber.Ctx) error {
	user := models.User{}
	if err := h.parseBody(c, &user); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
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
	if err := h.parseBody(c, &logReq); err != nil {
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

// User
func (h *Handler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(user)
}

// Wallets
func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Check if both fields are provided
	if len(wallet.Name) == 0 || len(wallet.Currency) == 0 {
		c.JSON(fiber.Map{"error": ErrWrongFormat.Error()})
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
	walletID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet, err := h.walletService.GetWalletByID(int64(walletID), userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(wallet)
}

func (h *Handler) UpdateWallet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	walletID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Check if at least one field is provided
	if len(wallet.Name) == 0 && len(wallet.Currency) == 0 {
		c.JSON(fiber.Map{"error": ErrWrongFormat.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet.ID = int64(walletID)
	err = h.walletService.UpdateWallet(&wallet, userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(wallet)
	return c.SendStatus(fiber.StatusOK)
}

// Categories
func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	category := models.Category{}
	if err := h.parseBody(c, &category); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category.OwnerID = userID
	id, err := h.categoryService.CreateCategory(&category)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(fiber.Map{"id": id})
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetCategories(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categories, err := h.categoryService.GetAllCategories(userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(categories)
}

func (h *Handler) GetCategoryByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categoryID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category, err := h.categoryService.GetCategoryByID(int64(categoryID), userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(category)
}

func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categoryID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category := models.Category{}
	if err := h.parseBody(c, &category); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category.ID = int64(categoryID)
	err = h.categoryService.UpdateCategory(&category, userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(category)
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categoryID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.categoryService.DeleteCategory(int64(categoryID), userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

// Operations
func (h *Handler) CreateOperation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := h.operationSerivce.CreateOperation(&operation, userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": id})
}

func (h *Handler) GetOperations(c *fiber.Ctx) error {
	walletID, err := strconv.Atoi(c.Queries()["wallet_id"])
	if err != nil {
		c.JSON(fiber.Map{"error": ErrWrongFormat.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID := c.Locals("user_id").(int64)
	operations, err := h.operationSerivce.GetOperationsByWalletID(int64(walletID), userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(operations)
}

func (h *Handler) GetOperationByID(c *fiber.Ctx) error {
	operationID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID := c.Locals("user_id").(int64)
	operation, err := h.operationSerivce.GetOperationByID(int64(operationID), userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(operation)
}

func (h *Handler) UpdateOperation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	operationID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	operation.ID = int64(operationID)
	err = h.operationSerivce.UpdateOperation(&operation, userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(operation)
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteOperation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	operationID, err := h.parseID(c)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.operationSerivce.DeleteOperation(int64(operationID), userID)
	if err != nil {
		c.JSON(fiber.Map{"error": err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) parseID(c *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return -1, ErrWrongFormat
	}
	return id, nil
}

func (h *Handler) parseBody(c *fiber.Ctx, out interface{}) error {
	// Parse the request body into the struct
	if err := c.BodyParser(out); err != nil {
		return ErrWrongFormat
	}

	// Validate required fields
	if err := h.validate.Struct(out); err != nil {
		return ErrWrongFormat
	}

	return nil
}
