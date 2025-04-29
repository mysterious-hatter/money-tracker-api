package handlers

import (
	"finances-backend/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateWallet(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	// Check if both fields are provIded
	if len(wallet.Name) == 0 || len(wallet.Currency) == 0 {
		return h.sendError(c, ErrWrongFormat, ErrNotAllRequiredFieldsProvided)
	}

	wallet.OwnerId = userId
	Id, err := h.walletService.CreateWallet(&wallet)
	if err != nil {
		return h.sendError(c, ErrCannotCreateWallet, err)
	}

	return h.send(c, fiber.StatusCreated, fiber.Map{"Id": Id})
}

func (h *Handler) GetWallets(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	wallets, err := h.walletService.GetAllWallets(userId)
	if err != nil {
		return h.sendError(c, ErrCannotGetWallets, err)
	}
	return h.send(c, fiber.StatusOK, wallets)
}

func (h *Handler) GetWalletById(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	walletId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	wallet, err := h.walletService.GetWalletById(int64(walletId), userId)
	if err != nil {
		return h.sendError(c, ErrCannotGetWallet, err)
	}

	return h.send(c, fiber.StatusOK, wallet)
}

func (h *Handler) UpdateWallet(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	walletId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	wallet := models.Wallet{}
	if err := h.parseBody(c, &wallet); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	// Check if at least one field is provIded
	if len(wallet.Name) == 0 && len(wallet.Currency) == 0 {
		return h.sendError(c, ErrWrongFormat, ErrNotAllRequiredFieldsProvided)
	}

	wallet.Id = int64(walletId)
	err = h.walletService.UpdateWallet(&wallet, userId)
	if err != nil {
		return h.sendError(c, ErrCannotUpdateWallet, err)
	}

	return h.send(c, fiber.StatusOK, wallet)
}
