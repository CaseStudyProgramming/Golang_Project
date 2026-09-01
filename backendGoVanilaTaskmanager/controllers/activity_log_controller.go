package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"taskmanager/models"
	"taskmanager/utils"
)

type ActivityLogController struct {
	service ActivityLogServiceInterface
}

// ActivityLogServiceInterface defines the interface for activity log service operations
type ActivityLogServiceInterface interface {
	LogActivity(userID int64, taskID *int64, action string, entityType string, entityID *int64, details string, ipAddress string, userAgent string) error
	GetUserActivityLogs(userID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error)
	GetTaskActivityLogs(taskID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error)
	GetActivityLogByID(id int64) (*models.ActivityLog, error)
}

func NewActivityLogController(service ActivityLogServiceInterface) *ActivityLogController {
	return &ActivityLogController{service: service}
}

// GET /activity-logs - Get current user's activity logs
func (c *ActivityLogController) GetUserActivityLogs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var page int
	var limit int

	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		page = 1
	} else {
		pageInt, err := strconv.Atoi(pageParam)
		if err != nil || pageInt < 1 {
			utils.ErrorResponse(w, http.StatusBadRequest, "page must be greater than 0")
			return
		}
		page = pageInt
	}

	limitParam := r.URL.Query().Get("limit")
	if limitParam == "" {
		limit = 10
	} else {
		limitInt, err := strconv.Atoi(limitParam)
		if err != nil || limitInt < 1 {
			utils.ErrorResponse(w, http.StatusBadRequest, "limit must be greater than 0")
			return
		}
		limit = limitInt
	}

	logs, meta, err := c.service.GetUserActivityLogs(userID, page, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "success", map[string]interface{}{
		"data": logs,
		"meta": meta,
	})
}

// GET /tasks/{id}/activity-logs - Get activity logs for a specific task
func (c *ActivityLogController) GetTaskActivityLogs(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var page int
	var limit int

	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		page = 1
	} else {
		pageInt, err := strconv.Atoi(pageParam)
		if err != nil || pageInt < 1 {
			utils.ErrorResponse(w, http.StatusBadRequest, "page must be greater than 0")
			return
		}
		page = pageInt
	}

	limitParam := r.URL.Query().Get("limit")
	if limitParam == "" {
		limit = 10
	} else {
		limitInt, err := strconv.Atoi(limitParam)
		if err != nil || limitInt < 1 {
			utils.ErrorResponse(w, http.StatusBadRequest, "limit must be greater than 0")
			return
		}
		limit = limitInt
	}

	logs, meta, err := c.service.GetTaskActivityLogs(taskID, page, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO: Verify user has access to this task's logs by checking if the task belongs to the user
	// This requires injecting TaskService or adding a method to check task ownership

	utils.SuccessResponse(w, http.StatusOK, "success", map[string]interface{}{
		"data": logs,
		"meta": meta,
	})
}

// GET /activity-logs/{id} - Get a specific activity log by ID
func (c *ActivityLogController) GetActivityLogByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log, err := c.service.GetActivityLogByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Activity log with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Verify the log belongs to the user
	if log.UserID != userID {
		utils.ErrorResponse(w, http.StatusForbidden, "You do not have permission to view this activity log")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Activity log found successfully", log)
}
