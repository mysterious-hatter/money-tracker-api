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
	ErrInternalServerError          error = errors.New("internal server error")
	ErrNotAllRequiredFieldsProvided error = errors.New("not all required fields provided")
	ErrWrongFormat                  error = errors.New("wrong format")
	ErrAuthFailed                   error = errors.New("authentication failed")
	// Users
	ErrCannotCreateUser error = errors.New("cannot create user")
	ErrCannotGetProfile error = errors.New("cannot get profile")
	// Wallets
	ErrCannotCreateWallet error = errors.New("cannot create wallet")
	ErrCannotGetWallet    error = errors.New("cannot get wallet")
	ErrCannotGetWallets   error = errors.New("cannot get wallets")
	ErrCannotUpdateWallet error = errors.New("cannot update wallet")
	ErrCannotDeleteWallet error = errors.New("cannot delete wallet")
	// Categories
	ErrCannotCreateCategory error = errors.New("cannot create category")
	ErrCannotGetCategory    error = errors.New("cannot get category")
	ErrCannotGetCategories  error = errors.New("cannot get categories")
	ErrCannotUpdateCategory error = errors.New("cannot update category")
	ErrCannotDeleteCategory error = errors.New("cannot delete category")
	// Operations
	ErrCannotCreateOperation error = errors.New("cannot create operation")
	ErrCannotGetOperation    error = errors.New("cannot get operation")
	ErrCannotGetOperations   error = errors.New("cannot get operations")
	ErrCannotUpdateOperation error = errors.New("cannot update operation")
	ErrCannotDeleteOperation error = errors.New("cannot delete operation")
	// Search
	ErrCannotSearch error = errors.New("cannot to search")
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
	searchService    services.SearchService
	validate         *validator.Validate
}

func NewHandler(
	as services.AuthService,
	us services.UserService,
	ws services.WalletService,
	cs services.CategoryService,
	ops services.OperationService,
	ss services.SearchService,
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
		searchService:    ss,
		validate:         validate,
	}
}

func (h *Handler) parseId(c *fiber.Ctx) (int, error) {
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

	// Check if there is useful information in the handlers error
	switch err {
	case ErrWrongFormat, ErrNotAllRequiredFieldsProvided:
		// Syntax errors in the request
		statusCode = fiber.StatusBadRequest // 400
	case ErrCannotCreateUser, ErrCannotCreateWallet, ErrCannotCreateCategory, ErrCannotCreateOperation:
		// Normally means, that the user is trying to create an already existing entity
		statusCode = fiber.StatusConflict // 409
	case ErrAuthFailed:
		// Authentication errors
		statusCode = fiber.StatusUnauthorized // 401
	default:
		// All other errors
		statusCode = fiber.StatusInternalServerError // 500
	}

	// Check if there is useful information in the services error
	switch fullErr {
	case services.ErrAccessDenied:
		statusCode = fiber.StatusForbidden // 403
	case services.ErrUserNotFound, services.ErrCategoryNotFound, services.ErrWalletNotFound, services.ErrOperationNotFound:
		statusCode = fiber.StatusNotFound // 404
	case services.ErrNoCategoriesFound, services.ErrNoWalletsFound, services.ErrNoOperationsFound:
		// Maybe these errors will be removed in the future
		statusCode = fiber.StatusNotFound
	}

	if err == nil || fullErr == nil {
		err = ErrInternalServerError
		fullErr = errors.New("an unexpected error occurred, please try again later")
	}

	return c.Status(statusCode).JSON(ErrorResponse{Error: err.Error(), Description: fullErr.Error()})
}
