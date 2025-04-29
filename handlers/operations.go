package handlers

import (
	"finances-backend/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateOperation(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	id, err := h.operationSerivce.CreateOperation(&operation, userID)
	if err != nil {
		return h.sendError(c, ErrCannotCreateOperation, err)
	}

	return h.send(c, fiber.StatusCreated, fiber.Map{"id": id})
}

func (h *Handler) GetOperations(c *fiber.Ctx) error {
	walletID, err := strconv.Atoi(c.Queries()["wallet_id"])
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	sinceParam := c.Queries()["since"]

	userID := c.Locals("userId").(int64)

	operations := []models.Operation{}
	if len(sinceParam) > 9 { // DD-MM-YYYY
		sinceDate, err := time.Parse("02-01-2006", sinceParam)
		if err != nil {
			return h.sendError(c, ErrWrongFormat, err)
		}
		operations, err = h.operationSerivce.GetOperationsSinceDateByWalletID(int64(walletID), userID, sinceDate)
	} else {
		operations, err = h.operationSerivce.GetOperationsByWalletID(int64(walletID), userID)
	}

	if err != nil {
		return h.sendError(c, ErrCannotGetOperations, err)
	}
	return h.send(c, fiber.StatusOK, operations)
}

func (h *Handler) GetOperationByID(c *fiber.Ctx) error {
	operationID, err := h.parseID(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	userID := c.Locals("userId").(int64)
	operation, err := h.operationSerivce.GetOperationByID(int64(operationID), userID)
	if err != nil {
		return h.sendError(c, ErrCannotGetOperation, err)
	}

	return h.send(c, fiber.StatusOK, operation)
}

func (h *Handler) UpdateOperation(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	operationID, err := h.parseID(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	operation.ID = int64(operationID)
	err = h.operationSerivce.UpdateOperation(&operation, userID)
	if err != nil {
		return h.sendError(c, ErrCannotUpdateOperation, err)
	}

	return h.send(c, fiber.StatusOK, operation)
}

func (h *Handler) DeleteOperation(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	operationID, err := h.parseID(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	err = h.operationSerivce.DeleteOperation(int64(operationID), userID)
	if err != nil {
		return h.sendError(c, ErrCannotDeleteOperation, err)
	}

	return h.send(c, fiber.StatusNoContent, nil)
}
