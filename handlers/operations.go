package handlers

import (
	"finances-backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateOperation(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	// Check if both fields are provIded
	if len(operation.Name) == 0 || len(operation.Place) == 0 {
		return h.sendError(c, ErrWrongFormat, ErrNotAllRequiredFieldsProvided)
	}

	id, err := h.operationSerivce.CreateOperation(&operation, userId)
	if err != nil {
		return h.sendError(c, ErrCannotCreateOperation, err)
	}

	return h.send(c, fiber.StatusCreated, fiber.Map{"id": id})
}

func (h *Handler) GetOperations(c *fiber.Ctx) error {
	walletId, err := strconv.Atoi(c.Queries()["walletId"])
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	sinceParam := c.Queries()["since"]
	sortBy := c.Queries()["sortBy"]

	userId := c.Locals("userId").(int64)

	var operations []models.Operation

	var sinceDate models.DateOnly

	if len(sinceParam) > 9 { // DD-MM-YYYY
		sinceDate, err = models.ParseDateOnly(sinceParam)
		if err != nil {
			return h.sendError(c, ErrWrongFormat, err)
		}
	}
	operations, err = h.operationSerivce.GetOperations(userId, int64(walletId), sinceDate, sortBy)

	if err != nil {
		return h.sendError(c, ErrCannotGetOperations, err)
	}
	return h.send(c, fiber.StatusOK, operations)
}

func (h *Handler) GetOperationById(c *fiber.Ctx) error {
	operationId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	userId := c.Locals("userId").(int64)
	operation, err := h.operationSerivce.GetOperationById(int64(operationId), userId)
	if err != nil {
		return h.sendError(c, ErrCannotGetOperation, err)
	}

	return h.send(c, fiber.StatusOK, operation)
}

func (h *Handler) UpdateOperation(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	operationId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	operation.Id = int64(operationId)
	err = h.operationSerivce.UpdateOperation(&operation, userId)
	if err != nil {
		return h.sendError(c, ErrCannotUpdateOperation, err)
	}

	return h.send(c, fiber.StatusOK, operation)
}

func (h *Handler) DeleteOperation(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	operationId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	err = h.operationSerivce.DeleteOperation(int64(operationId), userId)
	if err != nil {
		return h.sendError(c, ErrCannotDeleteOperation, err)
	}

	return h.send(c, fiber.StatusNoContent, nil)
}
