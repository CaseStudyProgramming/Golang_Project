package services

import (
	"errors"
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
		return nil, errors.New("Task not found or access denied")
	}

	subtask.TaskID = task.ID
	subtask.IsCompleted = false

	createdSubtask, err := s.subtaskModel.Create(subtask)
	if err != nil {
		return nil, err
	}

	// Update task progress after creating subtask
	if err := s.taskModel.UpdateProgress(taskID); err != nil {
		return nil, err
	}

	return createdSubtask, nil
}

func (s *SubtaskService) GetByTaskID(userID int64, taskID int64) ([]models.Subtask, error) {
	// Check if task exists and belongs to user
	_, err := s.taskModel.GetByID(userID, taskID)
	if err != nil {
		return nil, errors.New("Task not found or access denied")
	}

	subtasks, err := s.subtaskModel.GetByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	return subtasks, nil
}

func (s *SubtaskService) GetByID(userID int64, id int64) (*models.Subtask, error) {
	subtask, err := s.subtaskModel.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, subtask.TaskID)
	if err != nil {
		return nil, errors.New("Access denied to subtask")
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
		return nil, err
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, existingSubtask.TaskID)
	if err != nil {
		return nil, errors.New("Access denied to subtask")
	}

	subtask.ID = id
	subtask.TaskID = existingSubtask.TaskID

	if err := s.subtaskModel.Update(subtask); err != nil {
		return nil, err
	}

	// Update task progress after updating subtask
	if err := s.taskModel.UpdateProgress(existingSubtask.TaskID); err != nil {
		return nil, err
	}

	return subtask, nil
}

func (s *SubtaskService) Delete(userID int64, id int64) error {
	// Get existing subtask to verify access
	existingSubtask, err := s.subtaskModel.GetByID(id)
	if err != nil {
		return err
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, existingSubtask.TaskID)
	if err != nil {
		return errors.New("Access denied to subtask")
	}

	if err := s.subtaskModel.Delete(id); err != nil {
		return err
	}

	// Update task progress after deleting subtask
	if err := s.taskModel.UpdateProgress(existingSubtask.TaskID); err != nil {
		return err
	}

	return nil
}

func (s *SubtaskService) Toggle(userID int64, id int64) error {
	// Get existing subtask to verify access
	existingSubtask, err := s.subtaskModel.GetByID(id)
	if err != nil {
		return err
	}

	// Check if the parent task belongs to the user
	_, err = s.taskModel.GetByID(userID, existingSubtask.TaskID)
	if err != nil {
		return errors.New("Access denied to subtask")
	}

	if err := s.subtaskModel.Toggle(id); err != nil {
		return err
	}

	// Update task progress after toggling subtask
	if err := s.taskModel.UpdateProgress(existingSubtask.TaskID); err != nil {
		return err
	}

	return nil
}
