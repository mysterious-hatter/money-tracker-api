package handlers

import (
	"finances-backend/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)

	category := models.Category{}
	if err := h.parseBody(c, &category); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	category.OwnerId = userId
	id, err := h.categoryService.CreateCategory(&category)
	if err != nil {
		return h.sendError(c, ErrCannotCreateCategory, err)
	}

	return h.send(c, fiber.StatusCreated, fiber.Map{"id": id})
}

func (h *Handler) GetCategories(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	categories, err := h.categoryService.GetAllCategories(userId)
	if err != nil {
		return h.sendError(c, ErrCannotGetCategories, err)
	}

	return h.send(c, fiber.StatusOK, categories)
}

func (h *Handler) GetCategoryById(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	categoryId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	category, err := h.categoryService.GetCategoryById(int64(categoryId), userId)
	if err != nil {
		return h.sendError(c, ErrCannotGetCategory, err)
	}

	return h.send(c, fiber.StatusOK, category)
}

func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	categoryId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	category := models.Category{}
	if err := h.parseBody(c, &category); err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	category.Id = int64(categoryId)
	err = h.categoryService.UpdateCategory(&category, userId)
	if err != nil {
		return h.sendError(c, ErrCannotUpdateCategory, err)
	}

	return h.send(c, fiber.StatusOK, category)
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	categoryId, err := h.parseId(c)
	if err != nil {
		return h.sendError(c, ErrWrongFormat, err)
	}

	err = h.categoryService.DeleteCategory(int64(categoryId), userId)
	if err != nil {
		return h.sendError(c, ErrCannotDeleteCategory, err)
	}

	return h.send(c, fiber.StatusNoContent, nil)
}
