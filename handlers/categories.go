package handlers

import (
	"finances-backend/models"
	"finances-backend/services"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	category := models.Category{}
	if err := h.parseBody(c, &category); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category.OwnerID = userID
	id, err := h.categoryService.CreateCategory(&category)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotCreateCategory.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(fiber.Map{"id": id})
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetCategories(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categories, err := h.categoryService.GetAllCategories(userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotGetCategories.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(categories)
}

func (h *Handler) GetCategoryByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categoryID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category, err := h.categoryService.GetCategoryByID(int64(categoryID), userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotGetCategory.Error(), Description: err.Error()})
		if err == services.ErrCategoryNotFound {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if err == services.ErrAccessDenied {
			return c.SendStatus(fiber.StatusForbidden)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(category)
}

func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categoryID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category := models.Category{}
	if err := h.parseBody(c, &category); err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	category.ID = int64(categoryID)
	err = h.categoryService.UpdateCategory(&category, userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotUpdateCategory.Error(), Description: err.Error()})
		if err == services.ErrCategoryNotFound {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if err == services.ErrAccessDenied {
			return c.SendStatus(fiber.StatusForbidden)
		}
		// Handle other errors (e.g., database errors)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.JSON(category)
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	categoryID, err := h.parseID(c)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrWrongFormat.Error(), Description: err.Error()})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = h.categoryService.DeleteCategory(int64(categoryID), userID)
	if err != nil {
		c.JSON(ErrorResponse{Error: ErrCannotDeleteCategory.Error(), Description: err.Error()})
		if err == services.ErrCategoryNotFound {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if err == services.ErrAccessDenied {
			return c.SendStatus(fiber.StatusForbidden)
		}
		// Handle other errors (e.g., database errors)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
