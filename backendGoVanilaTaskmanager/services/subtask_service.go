package services

import (
	"errors"
	"fmt"
	"taskmanager/models"
)

type SubtaskService struct {
	subtaskModel models.SubtaskModelInterface
	taskModel    models.TaskModelInterface
}

func NewSubtaskService(subtaskModel models.SubtaskModelInterface, taskModel models.TaskModelInterface) *SubtaskService {
	return &SubtaskService{
		subtaskModel: subtaskModel,
		taskModel:    taskModel,
	}
}

func (s *SubtaskService) Create(userID int64, taskID int64, subtask *models.Subtask) (*models.Subtask, error) {
	// Validate subtask title
	if subtask.Title == "" {
		return nil, errors.New("Subtask title cannot be empty")
	}

	// Check if task exists and belongs to user
	task, err := s.taskModel.GetByID(userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found or access denied: %w", err)
	}

	subtask.TaskID = task.ID
	subtask.IsCompleted = false

	createdSubtask, err := s.subtaskModel.Create(subtask)
	if err != nil {
		return nil, fmt.Errorf("failed to create subtask: %w", err)
	}

	// Update task progress after creating subtask
	if err := s.taskModel.UpdateProgress(taskID); err != nil {
		return nil, fmt.Errorf("failed to update task progress: %w", err)
	}

	return createdSubtask, nil
}

func (s *SubtaskService) GetByTaskID(userID int64, taskID int64) ([]models.Subtask, error) {
	// Check if task exists and belongs to user
	_, err := s.taskModel.GetByID(userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found or access denied: %w", err)
	}

	subtasks, err := s.subtaskModel.GetByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subtasks for task %d: %w", taskID, err)
	}

	return subtasks, nil
}

func (s *SubtaskService) GetByID(userID int64, id int64) (*models.Subtask, error) {
	subtask, err := s.subtaskModel.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subtask %d: %w", id, err)
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, subtask.TaskID)
	if err != nil {
		return nil, fmt.Errorf("access denied to subtask %d: %w", id, err)
	}

	return subtask, nil
}

func (s *SubtaskService) Update(userID int64, id int64, subtask *models.Subtask) (*models.Subtask, error) {
	// Validate subtask title
	if subtask.Title == "" {
		return nil, errors.New("Subtask title cannot be empty")
	}

	// Get existing subtask to verify access
	existingSubtask, err := s.subtaskModel.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subtask %d for update: %w", id, err)
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, existingSubtask.TaskID)
	if err != nil {
		return nil, fmt.Errorf("access denied to subtask %d: %w", id, err)
	}

	subtask.ID = id
	subtask.TaskID = existingSubtask.TaskID

	if err := s.subtaskModel.Update(subtask); err != nil {
		return nil, fmt.Errorf("failed to update subtask %d: %w", id, err)
	}

	// Update task progress after updating subtask
	if err := s.taskModel.UpdateProgress(existingSubtask.TaskID); err != nil {
		return nil, fmt.Errorf("failed to update task progress: %w", err)
	}

	return subtask, nil
}

func (s *SubtaskService) Delete(userID int64, id int64) error {
	// Get existing subtask to verify access
	existingSubtask, err := s.subtaskModel.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get subtask %d for deletion: %w", id, err)
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, existingSubtask.TaskID)
	if err != nil {
		return fmt.Errorf("access denied to subtask %d: %w", id, err)
	}

	if err := s.subtaskModel.Delete(id); err != nil {
		return fmt.Errorf("failed to delete subtask %d: %w", id, err)
	}

	// Update task progress after deleting subtask
	if err := s.taskModel.UpdateProgress(existingSubtask.TaskID); err != nil {
		return fmt.Errorf("failed to update task progress: %w", err)
	}

	return nil
}

func (s *SubtaskService) Toggle(userID int64, id int64) error {
	// Get existing subtask to verify access
	existingSubtask, err := s.subtaskModel.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get subtask %d for toggle: %w", id, err)
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, existingSubtask.TaskID)
	if err != nil {
		return fmt.Errorf("access denied to subtask %d: %w", id, err)
	}

	if err := s.subtaskModel.Toggle(id); err != nil {
		return fmt.Errorf("failed to toggle subtask %d: %w", id, err)
	}

	// Update task progress after toggling subtask
	if err := s.taskModel.UpdateProgress(existingSubtask.TaskID); err != nil {
		return fmt.Errorf("failed to update task progress: %w", err)
	}

	return nil
}
