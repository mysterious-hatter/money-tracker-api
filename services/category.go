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
		return 0, err
	}
	return id, nil
}

func (cs *CategoryService) GetAllCategories(userId int64) ([]models.Category, error) {
	categories, err := cs.storage.GetAllCategories(userId)
	if err != nil {
		return nil, ErrNoCategoriesFound
	}
	return categories, nil
}

func (cs *CategoryService) GetCategoryById(categoryId int64, userId int64) (*models.Category, error) {
	category, err := cs.storage.GetCategoryById(categoryId)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// Check user's ownership of the category
	if err := checkOwnership(category.OwnerId, userId); err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) UpdateCategory(category *models.Category, userId int64) error {
	// Check user's ownership of the category and if it exists
	_, err := cs.GetCategoryById(category.Id, userId)
	if err != nil {
		return err
	}

	err = cs.storage.UpdateCategory(category)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CategoryService) DeleteCategory(categoryId int64, userId int64) error {
	// Check user's ownership of the category and if it exists
	_, err := cs.GetCategoryById(categoryId, userId)
	if err != nil {
		return err
	}

	err = cs.storage.DeleteCategory(categoryId)
	if err != nil {
		return err
	}

	return nil
}
