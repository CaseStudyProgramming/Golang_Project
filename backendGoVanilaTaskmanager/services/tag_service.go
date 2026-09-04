package services

import (
	"errors"
	"fmt"
	"taskmanager/models"
)

type TagService struct {
	model models.TagModelInterface
}

func NewTagService(model models.TagModelInterface) *TagService {
	return &TagService{model: model}
}

func (s *TagService) Create(userID int64, tag *models.Tag) (*models.Tag, error) {
	if tag.Name == "" {
		return nil, errors.New("Tag name tidak boleh kosong")
	}

	// Set default color if not provided
	if tag.ColorHex == "" {
		tag.ColorHex = "#10B981"
	}

	tag.UserID = userID
	tag, err := s.model.Create(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return tag, nil
}

func (s *TagService) GetAll(userID int64) ([]models.Tag, error) {
	tags, err := s.model.GetAll(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return tags, nil
}

func (s *TagService) GetByID(userID int64, id int64) (*models.Tag, error) {
	tag, err := s.model.GetByID(userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag %d: %w", id, err)
	}
	return tag, nil
}

func (s *TagService) Update(userID int64, id int64, tag *models.Tag) (*models.Tag, error) {
	tag.ID = id
	tag.UserID = userID

	// Check if tag exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag %d for update: %w", id, err)
	}

	if tag.Name == "" {
		return nil, errors.New("Tag name tidak boleh kosong")
	}

	if err := s.model.Update(userID, tag); err != nil {
		return nil, fmt.Errorf("failed to update tag %d: %w", id, err)
	}
	return tag, nil
}

func (s *TagService) Delete(userID int64, id int64) error {
	// Check if tag exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return fmt.Errorf("failed to get tag %d for deletion: %w", id, err)
	}
	if err := s.model.Delete(userID, id); err != nil {
		return fmt.Errorf("failed to delete tag %d: %w", id, err)
	}
	return nil
}

// Task-tag relationship operations
func (s *TagService) AddTagToTask(taskID int64, tagID int64) error {
	if err := s.model.AddTagToTask(taskID, tagID); err != nil {
		return fmt.Errorf("failed to add tag %d to task %d: %w", tagID, taskID, err)
	}
	return nil
}

func (s *TagService) RemoveTagFromTask(taskID int64, tagID int64) error {
	if err := s.model.RemoveTagFromTask(taskID, tagID); err != nil {
		return fmt.Errorf("failed to remove tag %d from task %d: %w", tagID, taskID, err)
	}
	return nil
}

func (s *TagService) GetTagsByTaskID(taskID int64) ([]models.Tag, error) {
	tags, err := s.model.GetTagsByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags for task %d: %w", taskID, err)
	}
	return tags, nil
}

func (s *TagService) GetTasksByTagID(tagID int64) ([]models.Task, error) {
	tasks, err := s.model.GetTasksByTagID(tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for tag %d: %w", tagID, err)
	}
	return tasks, nil
}
