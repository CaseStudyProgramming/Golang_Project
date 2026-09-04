package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"taskmanager/models"
	"taskmanager/utils"
)

type CategoryController struct {
	service CategoryServiceInterface
}

// CategoryServiceInterface defines the interface for category service operations
type CategoryServiceInterface interface {
	Create(userID int64, category *models.Category) (*models.Category, error)
	GetAll(userID int64) ([]models.Category, error)
	GetByID(userID int64, id int64) (*models.Category, error)
	Update(userID int64, id int64, category *models.Category) (*models.Category, error)
	Delete(userID int64, id int64) error
}

func NewCategoryController(service CategoryServiceInterface) *CategoryController {
	return &CategoryController{service: service}
}

// POST /categories
func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	createdCategory, err := c.service.Create(userID, &category)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusCreated, "Category created successfully", createdCategory, timezone)
}

// GET /categories
func (c *CategoryController) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	categories, err := c.service.GetAll(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Categories retrieved successfully", categories, timezone)
}

// GET /categories/{id}
func (c *CategoryController) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	category, err := c.service.GetByID(userID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Category with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Category found successfully", category, timezone)
}

// PUT /categories/{id}
func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedCategory, err := c.service.Update(userID, id, &category)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Category with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Category updated successfully", updatedCategory, timezone)
}

// DELETE /categories/{id}
func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Delete(userID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Category with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusAccepted, "Category deleted successfully", nil, timezone)
}
