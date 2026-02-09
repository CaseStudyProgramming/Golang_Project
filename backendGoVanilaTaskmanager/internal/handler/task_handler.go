package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"taskmanager/internal/entity"
	"taskmanager/internal/repository"
	response_test "taskmanager/pkg/response"
	"time"
)

type TaskHandler struct {
	Repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

// POST
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		response_test.ErrorResponse(w, http.StatusBadRequest, "Title tidak boleh kosong")
		return
	}

	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		response_test.ErrorResponse(w, http.StatusBadRequest, "Due date must be equal or greater than current date")
		return
	}

	task.Completed = false // default

	createdTask, err := h.Repo.Create(&task)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response_test.SuccessResponse(w, http.StatusCreated, "Task created successfully", createdTask)
}

// GET ALL DATA
// GET ALL DATA
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var completed *bool
	var page int
	var limit int
	var search string

	completedParam := r.URL.Query().Get("completed")
	if completedParam != "" {
		value, err := strconv.ParseBool(completedParam)
		if err != nil {
			response_test.ErrorResponse(w, http.StatusBadRequest, "completed must be true or false")
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
			response_test.ErrorResponse(w, http.StatusBadRequest, "page must be greater than 0")
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
			response_test.ErrorResponse(w, http.StatusBadRequest, "limit must be greater than 0")
			return
		}
		limit = limitInt
	}

	// Validate presence of `search` param: if provided it must not be empty/whitespace
	vals, hasSearch := r.URL.Query()["search"]
	if hasSearch {
		if len(vals) == 0 || strings.TrimSpace(vals[0]) == "" {
			response_test.ErrorResponse(w, http.StatusBadRequest, "search must not be empty")
			return
		}
		search = vals[0]
	}

	offset := (page - 1) * limit

	tasks, total, err := h.Repo.GetAll(completed, offset, limit, search)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Validate page
	totalPages := (total + limit - 1) / limit // ceil division
	if total > 0 && page > totalPages {
		response_test.ErrorResponse(w, http.StatusBadRequest, "page must be less than or equal to "+strconv.Itoa(totalPages))
		return
	}

	// If user searched and no results found, return 404
	if strings.TrimSpace(search) != "" && len(tasks) == 0 {
		response_test.ErrorResponse(w, http.StatusNotFound, "no tasks found")
		return
	}

	// Add pagination meta data
	meta := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total_data": total,
		"total_page": totalPages,
		"has_next":   page < totalPages,
		"has_prev":   page > 1,
	}

	response_test.SuccessResponse(w, http.StatusOK, "success", map[string]interface{}{
		"data": tasks,
		"meta": meta,
	})
}

// GET BY ID
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.Repo.GetByID(id)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	response_test.SuccessResponse(w, http.StatusOK, "Task found successfully", task)
}

// PUT
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	var task entity.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	task.ID = id
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		response_test.ErrorResponse(w, http.StatusBadRequest, "Due date must be equal or greater than current date")
		return
	}
	if err := h.Repo.Update(&task); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	response_test.SuccessResponse(w, http.StatusOK, "Task updated successfully", task)
}

// PATCH
func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	var task entity.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	task.ID = id
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		response_test.ErrorResponse(w, http.StatusBadRequest, "Due date must be equal or greater than current date")
		return
	}
	if err := h.Repo.Complete(id); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	response_test.SuccessResponse(w, http.StatusOK, "Task completed successfully", nil)
}

// DELETE
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	softDelete := true
	if err := h.Repo.DeleteTask(id, softDelete); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
	response_test.SuccessResponse(w, http.StatusAccepted, "Task deleted successfully", nil)
}

// Mark Task as Completed /tasks/{id}/complete
func (h *TaskHandler) MarkTaskAsCompleted(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/") : len(r.URL.Path)-len("/complete")]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Repo.Complete(id); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	response_test.SuccessResponse(w, http.StatusOK, "Task completed successfully", nil)
}

// Mark Task as Uncompleted /tasks/{id}/uncomplete
func (h *TaskHandler) MarkTaskAsUncompleted(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/") : len(r.URL.Path)-len("/uncomplete")]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Repo.MarkTaskAsUncompleted(id); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	response_test.SuccessResponse(w, http.StatusOK, "Task uncompleted successfully", nil)
}
