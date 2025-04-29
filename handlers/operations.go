package handlers

import (
	"finances-backend/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateOperation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	id, err := h.operationSerivce.CreateOperation(&operation, userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotCreateOperation.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": id})
}

func (h *Handler) GetOperations(c *fiber.Ctx) error {
	walletID, err := strconv.Atoi(c.Queries()["wallet_id"])
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	sinceParam := c.Queries()["since"]

	userID := c.Locals("user_id").(int64)

	operations := []models.Operation{}
	if len(sinceParam) > 9 { // DD-MM-YYYY
		sinceDate, err := time.Parse("02-01-2006", sinceParam)
		if err != nil {
			c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
			return c.SendStatus(fiber.StatusBadRequest)
		}
		operations, err = h.operationSerivce.GetOperationsSinceDateByWalletID(int64(walletID), userID, sinceDate)
	} else {
		operations, err = h.operationSerivce.GetOperationsByWalletID(int64(walletID), userID)
	}

	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotGetOperations.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(operations)
}

func (h *Handler) GetOperationByID(c *fiber.Ctx) error {
	operationID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID := c.Locals("user_id").(int64)
	operation, err := h.operationSerivce.GetOperationByID(int64(operationID), userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotGetOperation.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(operation)
}

func (h *Handler) UpdateOperation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	operationID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	operation := models.Operation{}
	if err := h.parseBody(c, &operation); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	operation.ID = int64(operationID)
	err = h.operationSerivce.UpdateOperation(&operation, userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotUpdateOperation.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(operation)
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteOperation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	operationID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.operationSerivce.DeleteOperation(int64(operationID), userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotDeleteOperation.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
