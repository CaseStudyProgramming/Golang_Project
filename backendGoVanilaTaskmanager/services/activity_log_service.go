package services

import (
	"errors"
	"fmt"
	"taskmanager/models"
)

type ActivityLogService struct {
	model models.ActivityLogModelInterface
}

func NewActivityLogService(model models.ActivityLogModelInterface) *ActivityLogService {
	return &ActivityLogService{model: model}
}

// LogActivity creates a new activity log entry
func (s *ActivityLogService) LogActivity(userID int64, taskID *int64, action string, entityType string, entityID *int64, details string, ipAddress string, userAgent string) error {
	log := &models.ActivityLog{
		UserID:     userID,
		TaskID:     taskID,
		Action:     models.ActionType(action),
		EntityType: models.EntityType(entityType),
		EntityID:   entityID,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	_, err := s.model.Create(log)
	return err
}

// GetUserActivityLogs retrieves activity logs for a specific user with pagination
func (s *ActivityLogService) GetUserActivityLogs(userID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
	if page < 1 {
		return nil, nil, errors.New("page must be greater than 0")
	}
	if limit < 1 {
		return nil, nil, errors.New("limit must be greater than 0")
	}

	offset := (page - 1) * limit
	logs, total, err := s.model.GetByUserID(userID, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (total + limit - 1) / limit
	if total > 0 && page > totalPages {
		return nil, nil, fmt.Errorf("page must be less than or equal to %d", totalPages)
	}

	meta := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total_data": total,
		"total_page": totalPages,
		"has_next":   page < totalPages,
		"has_prev":   page > 1,
	}

	return logs, meta, nil
}

// GetTaskActivityLogs retrieves activity logs for a specific task with pagination
func (s *ActivityLogService) GetTaskActivityLogs(taskID int64, page int, limit int) ([]models.ActivityLog, map[string]interface{}, error) {
	if page < 1 {
		return nil, nil, errors.New("page must be greater than 0")
	}
	if limit < 1 {
		return nil, nil, errors.New("limit must be greater than 0")
	}

	offset := (page - 1) * limit
	logs, total, err := s.model.GetByTaskID(taskID, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (total + limit - 1) / limit
	if total > 0 && page > totalPages {
		return nil, nil, fmt.Errorf("page must be less than or equal to %d", totalPages)
	}

	meta := map[string]interface{}{
		"page":       page,
		"limit":      limit,
		"total_data": total,
		"total_page": totalPages,
		"has_next":   page < totalPages,
		"has_prev":   page > 1,
	}

	return logs, meta, nil
}

// GetActivityLogByID retrieves a single activity log by ID
func (s *ActivityLogService) GetActivityLogByID(id int64) (*models.ActivityLog, error) {
	log, err := s.model.GetByID(id)
	if err != nil {
		return nil, err
	}
	return log, nil
}
