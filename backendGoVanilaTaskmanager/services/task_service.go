package services

import (
	"errors"
	"fmt"
	"strings"
	"taskmanager/models"
	"time"
)

type TaskService struct {
	model models.TaskModelInterface
}

func NewTaskService(model models.TaskModelInterface) *TaskService {
	return &TaskService{model: model}
}

func (s *TaskService) Create(task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, errors.New("Title tidak boleh kosong")
	}

	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		return nil, errors.New("Due date must be equal or greater than current date")
	}

	task.Completed = false
	return s.model.Create(task)
}

func (s *TaskService) GetAll(completed *bool, page, limit int, search string) ([]models.Task, map[string]interface{}, error) {
	offset := (page - 1) * limit
	tasks, total, err := s.model.GetAll(completed, offset, limit, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (total + limit - 1) / limit
	if total > 0 && page > totalPages {
		return nil, nil, fmt.Errorf("page must be less than or equal to %d", totalPages)
	}

	if strings.TrimSpace(search) != "" && len(tasks) == 0 {
		return nil, nil, errors.New("no tasks found")
	}

	meta := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total_data": total,
		"total_page": totalPages,
		"has_next":   page < totalPages,
		"has_prev":   page > 1,
	}

	return tasks, meta, nil
}

func (s *TaskService) GetByID(id int64) (*models.Task, error) {
	return s.model.GetByID(id)
}

func (s *TaskService) Update(id int64, task *models.Task) (*models.Task, error) {
	task.ID = id
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		return nil, errors.New("Due date must be equal or greater than current date")
	}

	// Check if task exists
	_, err := s.model.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.model.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Complete(id int64) error {
	return s.model.Complete(id)
}

func (s *TaskService) Delete(id int64) error {
	// Check if task exists
	_, err := s.model.GetByID(id)
	if err != nil {
		return err
	}
	return s.model.Delete(id, true)
}

func (s *TaskService) Restore(id int64) error {
	// The original controller check is:
	// task, err := h.Repo.GetByID(id)
	// if err != nil && err != sql.ErrNoRows { ... }
	// if task == nil || task.ID == 0 { ... }
	// We'll mimic this or handle it nicely in the controller.
	return s.model.RestoreTask(id)
}

func (s *TaskService) MarkAsCompleted(id int64) error {
	return s.model.MarkTaskAsCompleted(id)
}

func (s *TaskService) MarkAsUncompleted(id int64) error {
	return s.model.MarkTaskAsUncompleted(id)
}
