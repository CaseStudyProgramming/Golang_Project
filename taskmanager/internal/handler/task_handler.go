package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"taskmanager/internal/entity"
	"taskmanager/internal/repository"
	response_test "taskmanager/pkg/response"
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
		response_test.ErrorResponse(w, "Title tidak boleh kosong")
		return
	}

	task.Completed = false // default

	if err := h.Repo.Create(&task); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response_test.SuccessResponse(w, http.StatusCreated, "Task created successfully", task)
}

// GET ALL DATA
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
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
	if err := h.Repo.Update(&task); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	response_test.SuccessResponse(w, http.StatusOK, "Task updated successfully", task)
}

//P

// DELETE
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response_test.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.Repo.Delete(id); err != nil {
		response_test.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
	response_test.SuccessResponse(w, http.StatusAccepted, "Task deleted successfully", nil)
