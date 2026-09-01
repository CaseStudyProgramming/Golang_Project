package services

import (
	"errors"
	"fmt"
	"strings"
	"taskmanager/models"
	"time"
)

type TaskService struct {
	model    models.TaskModelInterface
	tagModel models.TagModelInterface
}

func NewTaskService(model models.TaskModelInterface, tagModel models.TagModelInterface) *TaskService {
	return &TaskService{model: model, tagModel: tagModel}
}

func (s *TaskService) Create(userID int64, task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, errors.New("Title tidak boleh kosong")
	}

	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		return nil, errors.New("Due date must be equal or greater than current date")
	}

	// Set default priority if not provided
	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}

	task.UserID = userID
	task.Completed = false
	return s.model.Create(task)
}

func (s *TaskService) GetAll(userID int64, completed *bool, page, limit int, search string, priority *models.Priority, categoryID *int64, sortBy string, sortOrder string) ([]models.Task, map[string]interface{}, error) {
	offset := (page - 1) * limit
	tasks, total, err := s.model.GetAll(userID, completed, offset, limit, search, priority, categoryID, sortBy, sortOrder)
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

	// Load tags and subtasks for each task
	for i := range tasks {
		if err := s.model.LoadTags(&tasks[i]); err != nil {
			return nil, nil, err
		}
		if err := s.model.LoadSubtasks(&tasks[i]); err != nil {
			return nil, nil, err
		}
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

func (s *TaskService) GetByID(userID int64, id int64) (*models.Task, error) {
	task, err := s.model.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	// Load tags for the task
	if err := s.model.LoadTags(task); err != nil {
		return nil, err
	}
	// Load subtasks for the task
	if err := s.model.LoadSubtasks(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Update(userID int64, id int64, task *models.Task) (*models.Task, error) {
	task.ID = id
	task.UserID = userID
	if task.DueDate != nil && task.DueDate.Before(time.Now()) {
		return nil, errors.New("Due date must be equal or greater than current date")
	}

	// Check if task exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return nil, err
	}

	if err := s.model.Update(userID, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Complete(userID int64, id int64) error {
	return s.model.Complete(userID, id)
}

func (s *TaskService) Delete(userID int64, id int64) error {
	// Check if task exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return err
	}
	return s.model.Delete(userID, id, true)
}

func (s *TaskService) Restore(userID int64, id int64) error {
	// The original controller check is:
	// task, err := h.Repo.GetByID(id)
	// if err != nil && err != sql.ErrNoRows { ... }
	// if task == nil || task.ID == 0 { ... }
	// We'll mimic this or handle it nicely in the controller.
	return s.model.RestoreTask(userID, id)
}

func (s *TaskService) MarkAsCompleted(userID int64, id int64) error {
	return s.model.MarkTaskAsCompleted(userID, id)
}

func (s *TaskService) MarkAsUncompleted(userID int64, id int64) error {
	return s.model.MarkTaskAsUncompleted(userID, id)
}

// Tag operations for tasks
func (s *TaskService) AddTagToTask(userID int64, taskID int64, tagID int64) error {
	// Check if task exists and belongs to user
	_, err := s.model.GetByID(userID, taskID)
	if err != nil {
		return err
	}
	// Check if tag exists and belongs to user
	_, err = s.tagModel.GetByID(userID, tagID)
	if err != nil {
		return err
	}
	return s.model.AddTagToTask(taskID, tagID)
}

func (s *TaskService) RemoveTagFromTask(userID int64, taskID int64, tagID int64) error {
	// Check if task exists and belongs to user
	_, err := s.model.GetByID(userID, taskID)
	if err != nil {
		return err
	}
	return s.model.RemoveTagFromTask(taskID, tagID)
}
