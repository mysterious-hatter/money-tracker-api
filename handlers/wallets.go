package handlers

import (
	"finances-backend/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Check if both fields are provided
	if len(wallet.Name) == 0 || len(wallet.Currency) == 0 {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: "Not all required fields provided"})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet.OwnerID = userID
	id, err := h.walletService.CreateWallet(&wallet)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotCreateWallet.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"id": id})
}

func (h *Handler) GetWallets(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	wallets, err := h.walletService.GetAllWallets(userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotGetWallets.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(wallets)
}

func (h *Handler) GetWalletByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	walletID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet, err := h.walletService.GetWalletByID(int64(walletID), userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotGetWallet.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(wallet)
}

func (h *Handler) UpdateWallet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	walletID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Check if at least one field is provided
	if len(wallet.Name) == 0 && len(wallet.Currency) == 0 {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: "Not all required fields provided"})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	wallet.ID = int64(walletID)
	err = h.walletService.UpdateWallet(&wallet, userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotUpdateWallet.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(wallet)
	return c.SendStatus(fiber.StatusOK)
}
