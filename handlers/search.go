package handlers

import (
	"finances-backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SearchOperations(c *fiber.Ctx) error {
	var walletId, categoryId int
	var date models.DateOnly
	var err error

	// UserId
	userId := c.Locals("userId").(int64)
	// Name
	name := c.Query("name")
	// WalletId
	walletIdParam := c.Query("walletId")
	if len(walletIdParam) == 0 {
		return h.sendError(c, ErrWrongFormat, ErrNotAllRequiredFieldsProvided)
	}
	walletId, err = strconv.Atoi(walletIdParam)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}
	// Date
	dateParam := c.Query("date")
	if len(dateParam) > 9 { // DD-MM-YYYY
		date, err = models.ParseDateOnly(dateParam)
		if err != nil {
			return h.sendError(c, ErrWrongFormat, err)
		}
	}
	// Place
	place := c.Query("place")
	// CategoryId
	categoryIdParam := c.Query("categoryId")
	if len(categoryIdParam) != 0 {
		categoryId, err = strconv.Atoi(c.Query("categoryId"))
		if err != nil {
			return h.sendError(c, ErrWrongFormat, err)
		}
	}
	// SortBy
	sortBy := c.Query("sortBy")

	operations, err := h.searchService.SearchOperations(userId, name, int64(walletId), date, place, int64(categoryId), sortBy)
	if err != nil {
		return h.sendError(c, ErrCannotSearch, err)
	}

	return h.send(c, fiber.StatusOK, operations)
}
