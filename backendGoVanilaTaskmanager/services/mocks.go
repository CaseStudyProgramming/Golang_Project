package services

import (
	"errors"
	"sync"
	"taskmanager/models"
)

// MockCategoryModel is a mock implementation of CategoryModelInterface for testing
type MockCategoryModel struct {
	categories   map[int64]*models.Category
	nextID       int64
	createError  error
	getAllError  error
	getByIDError error
	updateError  error
	deleteError  error
}

func NewMockCategoryModel() *MockCategoryModel {
	return &MockCategoryModel{
		categories: make(map[int64]*models.Category),
		nextID:     1,
	}
}

func (m *MockCategoryModel) Create(category *models.Category) (*models.Category, error) {
	if m.createError != nil {
		return nil, m.createError
	}

	category.ID = m.nextID
	m.nextID++
	m.categories[category.ID] = category
	return category, nil
}

func (m *MockCategoryModel) GetAll(userID int64) ([]models.Category, error) {
	if m.getAllError != nil {
		return nil, m.getAllError
	}

	var result []models.Category
	for _, cat := range m.categories {
		if cat.UserID == userID {
			result = append(result, *cat)
		}
	}
	return result, nil
}

func (m *MockCategoryModel) GetByID(userID int64, id int64) (*models.Category, error) {
	if m.getByIDError != nil {
		return nil, m.getByIDError
	}

	cat, exists := m.categories[id]
	if !exists || cat.UserID != userID {
		return nil, errors.New("category not found")
	}
	return cat, nil
}

func (m *MockCategoryModel) Update(userID int64, category *models.Category) error {
	if m.updateError != nil {
		return m.updateError
	}

	cat, exists := m.categories[category.ID]
	if !exists || cat.UserID != userID {
		return errors.New("category not found")
	}

	cat.Name = category.Name
	cat.ColorHex = category.ColorHex
	return nil
}

func (m *MockCategoryModel) Delete(userID int64, id int64) error {
	if m.deleteError != nil {
		return m.deleteError
	}

	cat, exists := m.categories[id]
	if !exists || cat.UserID != userID {
		return errors.New("category not found")
	}

	delete(m.categories, id)
	return nil
}

// MockTagModel is a mock implementation of TagModelInterface for testing
type MockTagModel struct {
	tags           map[int64]*models.Tag
	nextID         int64
	taskTags       map[int64][]int64 // taskID -> tagIDs
	tagTasks       map[int64][]int64 // tagID -> taskIDs
	createError    error
	getAllError    error
	getByIDError   error
	updateError    error
	deleteError    error
	addTagError    error
	removeTagError error
	getTagsError   error
	getTasksError  error
}

func NewMockTagModel() *MockTagModel {
	return &MockTagModel{
		tags:     make(map[int64]*models.Tag),
		taskTags: make(map[int64][]int64),
		tagTasks: make(map[int64][]int64),
		nextID:   1,
	}
}

func (m *MockTagModel) Create(tag *models.Tag) (*models.Tag, error) {
	if m.createError != nil {
		return nil, m.createError
	}

	tag.ID = m.nextID
	m.nextID++
	m.tags[tag.ID] = tag
	return tag, nil
}

func (m *MockTagModel) GetAll(userID int64) ([]models.Tag, error) {
	if m.getAllError != nil {
		return nil, m.getAllError
	}

	var result []models.Tag
	for _, tag := range m.tags {
		if tag.UserID == userID {
			result = append(result, *tag)
		}
	}
	return result, nil
}

func (m *MockTagModel) setGetByIDError(err error) {
	m.getByIDError = err
}

func (m *MockTagModel) GetByID(userID int64, id int64) (*models.Tag, error) {
	if m.getByIDError != nil {
		return nil, m.getByIDError
	}

	tag, exists := m.tags[id]
	if !exists || tag.UserID != userID {
		return nil, errors.New("tag not found")
	}
	return tag, nil
}

func (m *MockTagModel) Update(userID int64, tag *models.Tag) error {
	if m.updateError != nil {
		return m.updateError
	}

	existingTag, exists := m.tags[tag.ID]
	if !exists || existingTag.UserID != userID {
		return errors.New("tag not found")
	}

	existingTag.Name = tag.Name
	existingTag.ColorHex = tag.ColorHex
	return nil
}

func (m *MockTagModel) Delete(userID int64, id int64) error {
	if m.deleteError != nil {
		return m.deleteError
	}

	tag, exists := m.tags[id]
	if !exists || tag.UserID != userID {
		return errors.New("tag not found")
	}

	delete(m.tags, id)
	return nil
}

func (m *MockTagModel) AddTagToTask(taskID int64, tagID int64) error {
	if m.addTagError != nil {
		return m.addTagError
	}

	m.taskTags[taskID] = append(m.taskTags[taskID], tagID)
	m.tagTasks[tagID] = append(m.tagTasks[tagID], taskID)
	return nil
}

func (m *MockTagModel) RemoveTagFromTask(taskID int64, tagID int64) error {
	if m.removeTagError != nil {
		return m.removeTagError
	}

	// Remove from taskTags
	if tags, exists := m.taskTags[taskID]; exists {
		for i, tid := range tags {
			if tid == tagID {
				m.taskTags[taskID] = append(tags[:i], tags[i+1:]...)
				break
			}
		}
	}

	// Remove from tagTasks
	if tasks, exists := m.tagTasks[tagID]; exists {
		for i, tID := range tasks {
			if tID == taskID {
				m.tagTasks[tagID] = append(tasks[:i], tasks[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (m *MockTagModel) GetTagsByTaskID(taskID int64) ([]models.Tag, error) {
	if m.getTagsError != nil {
		return nil, m.getTagsError
	}

	var result []models.Tag
	if tagIDs, exists := m.taskTags[taskID]; exists {
		for _, tagID := range tagIDs {
			if tag, tagExists := m.tags[tagID]; tagExists {
				result = append(result, *tag)
			}
		}
	}
	return result, nil
}

func (m *MockTagModel) GetTasksByTagID(tagID int64) ([]models.Task, error) {
	if m.getTasksError != nil {
		return nil, m.getTasksError
	}

	// Return empty slice for simplicity in mock
	return []models.Task{}, nil
}

// MockActivityLogModel is a mock implementation of ActivityLogModelInterface for testing
type MockActivityLogModel struct {
	mu           sync.RWMutex
	logs         map[int64]*models.ActivityLog
	nextID       int64
	createError  error
	getUserError error
	getTaskError error
	getByIDError error
}

func NewMockActivityLogModel() *MockActivityLogModel {
	return &MockActivityLogModel{
		logs:   make(map[int64]*models.ActivityLog),
		nextID: 1,
	}
}

func (m *MockActivityLogModel) Create(log *models.ActivityLog) (*models.ActivityLog, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.createError != nil {
		return nil, m.createError
	}

	log.ID = m.nextID
	m.nextID++
	m.logs[log.ID] = log
	return log, nil
}

func (m *MockActivityLogModel) GetByUserID(userID int64, offset int, limit int) ([]models.ActivityLog, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.getUserError != nil {
		return nil, 0, m.getUserError
	}

	var result []models.ActivityLog
	for _, log := range m.logs {
		if log.UserID == userID {
			result = append(result, *log)
		}
	}

	// Apply pagination
	total := len(result)
	start := offset
	if start > len(result) {
		start = len(result)
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	if start >= end {
		return []models.ActivityLog{}, total, nil
	}

	return result[start:end], total, nil
}

func (m *MockActivityLogModel) GetByTaskID(taskID int64, offset int, limit int) ([]models.ActivityLog, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.getTaskError != nil {
		return nil, 0, m.getTaskError
	}

	var result []models.ActivityLog
	for _, log := range m.logs {
		if log.TaskID != nil && *log.TaskID == taskID {
			result = append(result, *log)
		}
	}

	// Apply pagination
	total := len(result)
	start := offset
	if start > len(result) {
		start = len(result)
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	if start >= end {
		return []models.ActivityLog{}, total, nil
	}

	return result[start:end], total, nil
}

func (m *MockActivityLogModel) GetByID(id int64) (*models.ActivityLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.getByIDError != nil {
		return nil, m.getByIDError
	}

	log, exists := m.logs[id]
	if !exists {
		return nil, errors.New("activity log not found")
	}
	return log, nil
}
