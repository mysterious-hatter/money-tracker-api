package services

import (
	"errors"

	"finances-backend/models"
	"finances-backend/storage"
)

var (
	ErrCategoryNotFound  error = errors.New("category not found")
	ErrNoCategoriesFound error = errors.New("no categories found")
)

type CategoryService struct {
	storage storage.Storage
}

func NewCategoryService(st storage.Storage) *CategoryService {
	cs := CategoryService{st}
	return &cs
}

func (cs *CategoryService) CreateCategory(category *models.Category) (int64, error) {
	id, err := cs.storage.CreateCategory(category)
	if err != nil {
		return 0, ErrUnableToCreate
	}
	return id, nil
}

func (cs *CategoryService) GetAllCategories(userID int64) ([]models.Category, error) {
	categories, err := cs.storage.GetAllCategories(userID)
	if err != nil {
		return nil, ErrNoCategoriesFound
	}
	return categories, nil
}

func (cs *CategoryService) GetCategoryByID(categoryID int64, userID int64) (*models.Category, error) {
	category, err := cs.storage.GetCategoryByID(categoryID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// Check user's ownership of the category
	if err := checkOwnership(category.OwnerID, userID); err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) UpdateCategory(category *models.Category, userID int64) error {
	// Check user's ownership of the category and if it exists
	_, err := cs.GetCategoryByID(category.ID, userID)
	if err != nil {
		return err
	}

	err = cs.storage.UpdateCategory(category)
	if err != nil {
		return ErrUnableToUpdate
	}

	return nil
}

func (cs *CategoryService) DeleteCategory(categoryID int64, userID int64) error {
	// Check user's ownership of the category and if it exists
	_, err := cs.GetCategoryByID(categoryID, userID)
	if err != nil {
		return err
	}

	err = cs.storage.DeleteCategory(categoryID)
	if err != nil {
		return ErrUnableToDelete
	}

	return nil
}
