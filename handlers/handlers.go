package handlers

import (
	"errors"
	"finances-backend/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	// Common
	ErrInternalServerError   error = errors.New("internal server error")
	ErrWrongFormat           error = errors.New("wrong format")
	ErrAuthFailed            error = errors.New("authentication failed")
	// Users
	ErrCannotCreateUser      error = errors.New("cannot create user")
	ErrCannotGetProfile      error = errors.New("cannot get profile")
	// Wallets
	ErrCannotCreateWallet    error = errors.New("cannot create wallet")
	ErrCannotGetWallet       error = errors.New("cannot get wallet")
	ErrCannotGetWallets      error = errors.New("cannot get wallets")
	ErrCannotUpdateWallet    error = errors.New("cannot update wallet")
	ErrCannotDeleteWallet    error = errors.New("cannot delete wallet")
	// Categories
	ErrCannotCreateCategory  error = errors.New("cannot create category")
	ErrCannotGetCategory     error = errors.New("cannot get category")
	ErrCannotGetCategories   error = errors.New("cannot get categories")
	ErrCannotUpdateCategory  error = errors.New("cannot update category")
	ErrCannotDeleteCategory  error = errors.New("cannot delete category")
	// Operations
	ErrCannotCreateOperation error = errors.New("cannot create operation")
	ErrCannotGetOperation    error = errors.New("cannot get operation")
	ErrCannotGetOperations   error = errors.New("cannot get operations")
	ErrCannotUpdateOperation error = errors.New("cannot update operation")
	ErrCannotDeleteOperation error = errors.New("cannot delete operation")
)

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

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
	validate := validator.New()
	err := validate.RegisterValidation("nonzero", func(fl validator.FieldLevel) bool {
		return fl.Field().Float() != 0
	})
	if err != nil {
		panic("failed to register custom validation: " + err.Error())
	}
	
	return &Handler{
		authService:      as,
		userService:      us,
		walletService:    ws,
		categoryService:  cs,
		operationSerivce: ops,
		validate:         validate,
	}
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
		return err
	}

	// Validate required fields
	if err := h.validate.Struct(out); err != nil {
		return err
	}

	return nil
}

func (h *Handler) send(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(data)
}

func (h *Handler) sendError(c *fiber.Ctx, err error, fullErr error) error {
	var statusCode int

	switch {
	case err == ErrWrongFormat:
		statusCode = fiber.StatusBadRequest
	case err == ErrAuthFailed:
		statusCode = fiber.StatusUnauthorized
	case err == ErrCannotGetProfile:
		statusCode = fiber.StatusInternalServerError
	case fullErr == services.ErrAccessDenied:
		statusCode = fiber.StatusForbidden
	case fullErr == services.ErrCategoryNotFound, fullErr == services.ErrWalletNotFound, fullErr == services.ErrOperationNotFound:
		statusCode = fiber.StatusNotFound
	// Maybe these errors will be removed in the future
	case fullErr == services.ErrNoCategoriesFound, fullErr == services.ErrNoWalletsFound, fullErr == services.ErrNoOperationsFound:
		statusCode = fiber.StatusNotFound
	default:
		statusCode = fiber.StatusInternalServerError
	}

	if err == nil || fullErr == nil {
		err = ErrInternalServerError
		fullErr = errors.New("an unexpected error occurred, please try again later")
	}

	return c.Status(statusCode).JSON(ErrorResponse{Error: err.Error(), Description: fullErr.Error()})
}