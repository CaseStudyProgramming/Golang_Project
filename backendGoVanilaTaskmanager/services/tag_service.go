package services

import (
	"errors"
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
	return s.model.Create(tag)
}

func (s *TagService) GetAll(userID int64) ([]models.Tag, error) {
	return s.model.GetAll(userID)
}

func (s *TagService) GetByID(userID int64, id int64) (*models.Tag, error) {
	return s.model.GetByID(userID, id)
}

func (s *TagService) Update(userID int64, id int64, tag *models.Tag) (*models.Tag, error) {
	tag.ID = id
	tag.UserID = userID

	// Check if tag exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return nil, err
	}

	if tag.Name == "" {
		return nil, errors.New("Tag name tidak boleh kosong")
	}

	if err := s.model.Update(userID, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) Delete(userID int64, id int64) error {
	// Check if tag exists
	_, err := s.model.GetByID(userID, id)
	if err != nil {
		return err
	}
	return s.model.Delete(userID, id)
}

// Task-tag relationship operations
func (s *TagService) AddTagToTask(taskID int64, tagID int64) error {
	return s.model.AddTagToTask(taskID, tagID)
}

func (s *TagService) RemoveTagFromTask(taskID int64, tagID int64) error {
	return s.model.RemoveTagFromTask(taskID, tagID)
}

func (s *TagService) GetTagsByTaskID(taskID int64) ([]models.Tag, error) {
	return s.model.GetTagsByTaskID(taskID)
}

func (s *TagService) GetTasksByTagID(tagID int64) ([]models.Task, error) {
	return s.model.GetTasksByTagID(tagID)
}
