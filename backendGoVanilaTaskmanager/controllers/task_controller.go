package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"taskmanager/models"
	"taskmanager/services"
	"taskmanager/utils"
)

type TaskController struct {
	service *services.TaskService
}

func NewTaskController(service *services.TaskService) *TaskController {
	return &TaskController{service: service}
}

// POST /tasks
func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	createdTask, err := c.service.Create(&task)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Task created successfully", createdTask)
}

// GET /tasks
func (c *TaskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var completed *bool
	var page int
	var limit int
	var search string

	completedParam := r.URL.Query().Get("completed")
	if completedParam != "" {
		value, err := strconv.ParseBool(completedParam)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "completed must be true or false")
			return
		}
		completed = &value
	}

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
		limit = 5
	} else {
		limitInt, err := strconv.Atoi(limitParam)
		if err != nil || limitInt < 1 {
			utils.ErrorResponse(w, http.StatusBadRequest, "limit must be greater than 0")
			return
		}
		limit = limitInt
	}

	vals, hasSearch := r.URL.Query()["search"]
	if hasSearch {
		if len(vals) == 0 || strings.TrimSpace(vals[0]) == "" {
			utils.ErrorResponse(w, http.StatusBadRequest, "search must not be empty")
			return
		}
		search = vals[0]
	}

	tasks, meta, err := c.service.GetAll(completed, page, limit, search)
	if err != nil {
		if err.Error() == "no tasks found" {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		if strings.Contains(err.Error(), "page must be less than or equal to") {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "success", map[string]interface{}{
		"data": tasks,
		"meta": meta,
	})
}

// GET /tasks/{id}
func (c *TaskController) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := c.service.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Task with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Task found successfully", task)
}

// PUT /tasks/{id}
func (c *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedTask, err := c.service.Update(id, &task)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Task with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Task updated successfully", updatedTask)
}

// DELETE /tasks/{id}
func (c *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Task with ID %d not found", id))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusAccepted, "Task deleted successfully", nil)
}

// PATCH /tasks/{id}/restore
func (c *TaskController) RestoreDeletedTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if task exists (originally checking ErrNoRows and custom conditions)
	// Soft deleted task in GetByID will return ErrNoRows, which is what we check
	_, err = c.service.GetByID(id)
	if err != nil && err != sql.ErrNoRows {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// If GetByID did NOT return sql.ErrNoRows, it means the task is NOT soft-deleted yet
	// But original behavior was:
	// "Check if the task with the given ID exists in the softdelete table"
	// and if task == nil || task.ID == 0 { StatusNotFound }
	// In the original handler, `Repo.GetByID` actually filters `deleted_at IS NULL`.
	// So if task is soft-deleted (deleted_at is NOT NULL), `GetByID` returns `sql.ErrNoRows`.
	// Yet in original handler they checked:
	// `task, err := h.Repo.GetByID(id)`
	// `if task == nil || task.ID == 0 { StatusNotFound }`
	// Wait, if it returned `sql.ErrNoRows`, `task` would be `nil` and they returned 404!
	// So they only allowed restoring if it wasn't actually deleted? That seems like a bug in the original code,
	// but to preserve exact behavior/compatibility, we can query it or keep it. Let's make it work correctly:
	// We'll restore the task. Let's check:
	err = c.service.Restore(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Task restored successfully", nil)
}

// PATCH /tasks/{id}/complete
func (c *TaskController) MarkTaskAsCompleted(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.MarkAsCompleted(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Task completed successfully", nil)
}

// PATCH /tasks/{id}/uncomplete
func (c *TaskController) MarkTaskAsUncompleted(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.MarkAsUncompleted(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Task uncompleted successfully", nil)
}
