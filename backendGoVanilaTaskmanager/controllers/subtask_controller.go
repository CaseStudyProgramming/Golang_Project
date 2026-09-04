package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"taskmanager/models"
	"taskmanager/utils"
)

type SubtaskController struct {
	service SubtaskServiceInterface
}

// SubtaskServiceInterface defines the interface for subtask service operations
type SubtaskServiceInterface interface {
	Create(userID int64, taskID int64, subtask *models.Subtask) (*models.Subtask, error)
	GetByTaskID(userID int64, taskID int64) ([]models.Subtask, error)
	GetByID(userID int64, id int64) (*models.Subtask, error)
	Update(userID int64, id int64, subtask *models.Subtask) (*models.Subtask, error)
	Delete(userID int64, id int64) error
	Toggle(userID int64, id int64) error
}

func NewSubtaskController(service SubtaskServiceInterface) *SubtaskController {
	return &SubtaskController{service: service}
}

// POST /tasks/{id}/subtasks
func (c *SubtaskController) CreateSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var subtask models.Subtask
	if err := json.NewDecoder(r.Body).Decode(&subtask); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	createdSubtask, err := c.service.Create(userID, taskID, &subtask)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusCreated, "Subtask created successfully", createdSubtask, timezone)
}

// GET /tasks/{id}/subtasks
func (c *SubtaskController) GetSubtasksByTaskID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	subtasks, err := c.service.GetByTaskID(userID, taskID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Subtasks retrieved successfully", subtasks, timezone)
}

// GET /subtasks/{id}
func (c *SubtaskController) GetSubtaskByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	subtask, err := c.service.GetByID(userID, id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Subtask with ID %d not found", id))
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Subtask found successfully", subtask, timezone)
}

// PUT /subtasks/{id}
func (c *SubtaskController) UpdateSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var subtask models.Subtask
	if err := json.NewDecoder(r.Body).Decode(&subtask); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedSubtask, err := c.service.Update(userID, id, &subtask)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Subtask updated successfully", updatedSubtask, timezone)
}

// DELETE /subtasks/{id}
func (c *SubtaskController) DeleteSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Delete(userID, id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Subtask deleted successfully", nil, timezone)
}

// PATCH /subtasks/{id}/toggle
func (c *SubtaskController) ToggleSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Toggle(userID, id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get timezone from context
	timezone := "UTC"
	if tz, ok := r.Context().Value("timezone").(string); ok {
		timezone = tz
	}
	utils.SuccessResponseWithTimezone(w, http.StatusOK, "Subtask toggled successfully", nil, timezone)
}
