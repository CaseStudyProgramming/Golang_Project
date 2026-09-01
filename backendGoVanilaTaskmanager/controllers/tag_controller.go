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

type TagController struct {
	service TagServiceInterface
}

// TagServiceInterface defines the interface for tag service operations
type TagServiceInterface interface {
	Create(userID int64, tag *models.Tag) (*models.Tag, error)
	GetAll(userID int64) ([]models.Tag, error)
	GetByID(userID int64, id int64) (*models.Tag, error)
	Update(userID int64, id int64, tag *models.Tag) (*models.Tag, error)
	Delete(userID int64, id int64) error
	GetTagsByTaskID(taskID int64) ([]models.Tag, error)
	GetTasksByTagID(tagID int64) ([]models.Task, error)
}

func NewTagController(service TagServiceInterface) *TagController {
	return &TagController{service: service}
}

// POST /tags
func (c *TagController) CreateTag(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	createdTag, err := c.service.Create(userID, &tag)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Tag created successfully", createdTag)
}

// GET /tags
func (c *TagController) GetAllTags(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	tags, err := c.service.GetAll(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Tags retrieved successfully", tags)
}

// GET /tags/{id}
func (c *TagController) GetTagByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	tag, err := c.service.GetByID(userID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Tag with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Tag found successfully", tag)
}

// PUT /tags/{id}
func (c *TagController) UpdateTag(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedTag, err := c.service.Update(userID, id, &tag)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Tag with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Tag updated successfully", updatedTag)
}

// DELETE /tags/{id}
func (c *TagController) DeleteTag(w http.ResponseWriter, r *http.Request) {
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
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Tag with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusAccepted, "Tag deleted successfully", nil)
}

// GET /tasks/{id}/tags
func (c *TagController) GetTagsByTaskID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	tags, err := c.service.GetTagsByTaskID(taskID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Tags retrieved successfully", tags)
}

// GET /tags/{id}/tasks
func (c *TagController) GetTasksByTagID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	tagID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	tasks, err := c.service.GetTasksByTagID(tagID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Tasks retrieved successfully", tasks)
}
