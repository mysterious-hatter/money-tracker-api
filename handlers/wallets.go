package handlers

import (
	"finances-backend/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	// Check if both fields are provided
	if len(wallet.Name) == 0 || len(wallet.Currency) == 0 {
		return h.sendError(c, ErrWrongFormat, ErrNotAllRequiredFieldsProvided)
	}

	wallet.OwnerID = userID
	id, err := h.walletService.CreateWallet(&wallet)
	if err != nil {
		return h.sendError(c, ErrCannotCreateWallet, err)
	}

	return h.send(c, fiber.StatusCreated, fiber.Map{"id": id})
}

func (h *Handler) GetWallets(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	wallets, err := h.walletService.GetAllWallets(userID)
	if err != nil {
		return h.sendError(c, ErrCannotGetWallets, err)
	}
	return h.send(c, fiber.StatusOK, wallets)
}

func (h *Handler) GetWalletByID(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	walletID, err := h.parseID(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	wallet, err := h.walletService.GetWalletByID(int64(walletID), userID)
	if err != nil {
		return h.sendError(c, ErrCannotGetWallet, err)
	}

	return h.send(c, fiber.StatusOK, wallet)
}

func (h *Handler) UpdateWallet(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	walletID, err := h.parseID(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	// Check if at least one field is provided
	if len(wallet.Name) == 0 && len(wallet.Currency) == 0 {
		return h.sendError(c, ErrWrongFormat, ErrNotAllRequiredFieldsProvided)
	}

	wallet.ID = int64(walletID)
	err = h.walletService.UpdateWallet(&wallet, userID)
	if err != nil {
		return h.sendError(c, ErrCannotUpdateWallet, err)
	}

	return h.send(c, fiber.StatusOK, wallet)
}
