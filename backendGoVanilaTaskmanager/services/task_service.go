package services

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"taskmanager/models"
	"time"
)

type TaskService struct {
	model           models.TaskModelInterface
	tagModel        models.TagModelInterface
	activityService ActivityLogServiceInterface
}

// ActivityLogServiceInterface defines the interface for activity logging
type ActivityLogServiceInterface interface {
	LogActivity(userID int64, taskID *int64, action string, entityType string, entityID *int64, details string, ipAddress string, userAgent string) error
}

func NewTaskService(model models.TaskModelInterface, tagModel models.TagModelInterface, activityService ActivityLogServiceInterface) *TaskService {
	return &TaskService{model: model, tagModel: tagModel, activityService: activityService}
}

func (s *TaskService) Create(userID int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
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
	createdTask, err := s.model.Create(task)
	if err != nil {
		return nil, err
	}

	// Log activity
	taskID := createdTask.ID
	if s.activityService != nil {
		details := fmt.Sprintf("Created task: %s", task.Title)
		_ = s.activityService.LogActivity(userID, &taskID, string(models.ActionCreate), string(models.EntityTask), &taskID, details, ipAddress, userAgent)
	}

	return createdTask, nil
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

func (s *TaskService) Update(userID int64, id int64, task *models.Task, ipAddress string, userAgent string) (*models.Task, error) {
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

	// Log activity
	if s.activityService != nil {
		details := fmt.Sprintf("Updated task: %s", task.Title)
		_ = s.activityService.LogActivity(userID, &id, string(models.ActionUpdate), string(models.EntityTask), &id, details, ipAddress, userAgent)
	}

	return task, nil
}

func (s *TaskService) Complete(userID int64, id int64) error {
	return s.model.Complete(userID, id)
}

func (s *TaskService) Delete(userID int64, id int64, ipAddress string, userAgent string) error {
	// Check if task exists
	task, err := s.model.GetByID(userID, id)
	if err != nil {
		return err
	}

	if err := s.model.Delete(userID, id, true); err != nil {
		return err
	}

	// Log activity
	if s.activityService != nil {
		details := fmt.Sprintf("Deleted task: %s", task.Title)
		_ = s.activityService.LogActivity(userID, &id, string(models.ActionDelete), string(models.EntityTask), &id, details, ipAddress, userAgent)
	}

	return nil
}

func (s *TaskService) Restore(userID int64, id int64) error {
	// The original controller check is:
	// task, err := h.Repo.GetByID(id)
	// if err != nil && err != sql.ErrNoRows { ... }
	// if task == nil || task.ID == 0 { ... }
	// We'll mimic this or handle it nicely in the controller.
	return s.model.RestoreTask(userID, id)
}

func (s *TaskService) MarkAsCompleted(userID int64, id int64, ipAddress string, userAgent string) error {
	// Get task for logging
	task, err := s.model.GetByID(userID, id)
	if err != nil {
		return err
	}

	if err := s.model.MarkTaskAsCompleted(userID, id); err != nil {
		return err
	}

	// Log activity
	if s.activityService != nil {
		details := fmt.Sprintf("Marked task as completed: %s", task.Title)
		_ = s.activityService.LogActivity(userID, &id, string(models.ActionComplete), string(models.EntityTask), &id, details, ipAddress, userAgent)
	}

	return nil
}

func (s *TaskService) MarkAsUncompleted(userID int64, id int64, ipAddress string, userAgent string) error {
	// Get task for logging
	task, err := s.model.GetByID(userID, id)
	if err != nil {
		return err
	}

	if err := s.model.MarkTaskAsUncompleted(userID, id); err != nil {
		return err
	}

	// Log activity
	if s.activityService != nil {
		details := fmt.Sprintf("Marked task as uncompleted: %s", task.Title)
		_ = s.activityService.LogActivity(userID, &id, string(models.ActionUncomplete), string(models.EntityTask), &id, details, ipAddress, userAgent)
	}

	return nil
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

func (s *TaskService) GetAnalyticsSummary(userID int64) (map[string]interface{}, error) {
	return s.model.GetAnalyticsSummary(userID)
}

func (s *TaskService) BulkDelete(userID int64, taskIDs []int64, ipAddress string, userAgent string) error {
	if len(taskIDs) == 0 {
		return nil
	}

	// Log activity for each task
	if s.activityService != nil {
		for _, taskID := range taskIDs {
			task, err := s.model.GetByID(userID, taskID)
			if err == nil {
				details := fmt.Sprintf("Bulk deleted task: %s", task.Title)
				_ = s.activityService.LogActivity(userID, &taskID, string(models.ActionDelete), string(models.EntityTask), &taskID, details, ipAddress, userAgent)
			}
		}
	}

	return s.model.BulkDelete(userID, taskIDs)
}

func (s *TaskService) BulkComplete(userID int64, taskIDs []int64, ipAddress string, userAgent string) error {
	if len(taskIDs) == 0 {
		return nil
	}

	// Log activity for each task
	if s.activityService != nil {
		for _, taskID := range taskIDs {
			task, err := s.model.GetByID(userID, taskID)
			if err == nil {
				details := fmt.Sprintf("Bulk completed task: %s", task.Title)
				_ = s.activityService.LogActivity(userID, &taskID, string(models.ActionComplete), string(models.EntityTask), &taskID, details, ipAddress, userAgent)
			}
		}
	}

	return s.model.BulkComplete(userID, taskIDs)
}

func (s *TaskService) ExportToCSV(userID int64, completed *bool, search string, priority *models.Priority, categoryID *int64) ([]byte, error) {
	// Get all tasks without pagination for export
	tasks, _, err := s.model.GetAll(userID, completed, 0, 0, search, priority, categoryID, "", "")
	if err != nil {
		return nil, err
	}

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	header := []string{"ID", "Title", "SubTitle", "Description", "Completed", "DueDate", "Priority", "CategoryID", "CreatedAt", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write task data
	for _, task := range tasks {
		dueDate := ""
		if task.DueDate != nil {
			dueDate = task.DueDate.Format("2006-01-02 15:04:05")
		}

		categoryIDVal := ""
		if task.CategoryID != nil {
			categoryIDVal = fmt.Sprintf("%d", *task.CategoryID)
		}

		record := []string{
			fmt.Sprintf("%d", task.ID),
			task.Title,
			task.SubTitle,
			task.Description,
			fmt.Sprintf("%t", task.Completed),
			dueDate,
			string(task.Priority),
			categoryIDVal,
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			task.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
