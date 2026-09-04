package services

import (
	"errors"
	"fmt"
	"taskmanager/models"
)

type CategoryService struct {
	model models.CategoryModelInterface
}

func NewCategoryService(model models.CategoryModelInterface) *CategoryService {
	return &CategoryService{model: model}
}

func (s *CategoryService) Create(userID int64, category *models.Category) (*models.Category, error) {
	if category.Name == "" {
		return nil, errors.New("Category name tidak boleh kosong")
	}

	// Set default color if not provided
	if category.ColorHex == "" {
		category.ColorHex = "#3B82F6"
	}

	category.UserID = userID
	return s.model.Create(category)
}

func (s *CategoryService) GetAll(userID int64) ([]models.Category, error) {
	return s.model.GetAll(userID)
}

func (s *CategoryService) GetByID(userID int64, id int64) (*models.Category, error) {
	return s.model.GetByID(userID, id)
}

func (s *CategoryService) Update(userID int64, id int64, category *models.Category) (*models.Category, error) {
	category.ID = id
	category.UserID = userID

	// Check if category exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category %d for update: %w", id, err)
	}

	if category.Name == "" {
		return nil, errors.New("Category name tidak boleh kosong")
	}

	if err := s.model.Update(userID, category); err != nil {
		return nil, fmt.Errorf("failed to update category %d: %w", id, err)
	}
	return category, nil
}

func (s *CategoryService) Delete(userID int64, id int64) error {
	// Check if category exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return fmt.Errorf("failed to get category %d for deletion: %w", id, err)
	}
	return s.model.Delete(userID, id)
}
