package models

import (
	"database/sql"
	"time"
)

// ActionType represents the type of action performed
type ActionType string

const (
	ActionCreate     ActionType = "CREATE"
	ActionUpdate     ActionType = "UPDATE"
	ActionDelete     ActionType = "DELETE"
	ActionRestore    ActionType = "RESTORE"
	ActionComplete   ActionType = "COMPLETE"
	ActionUncomplete ActionType = "UNCOMPLETE"
)

// EntityType represents the type of entity affected
type EntityType string

const (
	EntityTask     EntityType = "TASK"
	EntitySubtask  EntityType = "SUBTASK"
	EntityCategory EntityType = "CATEGORY"
	EntityTag      EntityType = "TAG"
)

// ActivityLog represents an activity log entry
type ActivityLog struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	TaskID     *int64     `json:"task_id,omitempty"`
	Action     ActionType `json:"action"`
	EntityType EntityType `json:"entity_type"`
	EntityID   *int64     `json:"entity_id,omitempty"`
	Details    string     `json:"details,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
	UserAgent  string     `json:"user_agent,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ActivityLogModel handles database operations for activity logs
type ActivityLogModel struct {
	DB *sql.DB
}

// ActivityLogModelInterface defines the interface for activity log model operations
type ActivityLogModelInterface interface {
	Create(log *ActivityLog) (*ActivityLog, error)
	GetByUserID(userID int64, offset, limit int) ([]ActivityLog, int, error)
	GetByTaskID(taskID int64, offset, limit int) ([]ActivityLog, int, error)
	GetByID(id int64) (*ActivityLog, error)
}

func NewActivityLogModel(db *sql.DB) *ActivityLogModel {
	return &ActivityLogModel{DB: db}
}

func (m *ActivityLogModel) Create(log *ActivityLog) (*ActivityLog, error) {
	query := `INSERT INTO activity_logs (user_id, task_id, action, entity_type, entity_id, details, ip_address, user_agent) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	          RETURNING id, user_id, task_id, action, entity_type, entity_id, details, ip_address, user_agent, created_at`
	err := m.DB.QueryRow(query, log.UserID, log.TaskID, log.Action, log.EntityType, log.EntityID, log.Details, log.IPAddress, log.UserAgent).
		Scan(&log.ID, &log.UserID, &log.TaskID, &log.Action, &log.EntityType, &log.EntityID, &log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt)
	if err != nil {
		return nil, err
	}
	return log, nil
}

func (m *ActivityLogModel) GetByUserID(userID int64, offset int, limit int) ([]ActivityLog, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM activity_logs WHERE user_id = $1`
	var total int
	if err := m.DB.QueryRow(countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := `SELECT id, user_id, task_id, action, entity_type, entity_id, details, ip_address, user_agent, created_at 
	          FROM activity_logs 
	          WHERE user_id = $1 
	          ORDER BY created_at DESC 
	          LIMIT $2 OFFSET $3`

	rows, err := m.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	logs := make([]ActivityLog, 0)
	for rows.Next() {
		var log ActivityLog
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.TaskID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.Details,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (m *ActivityLogModel) GetByTaskID(taskID int64, offset int, limit int) ([]ActivityLog, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM activity_logs WHERE task_id = $1`
	var total int
	if err := m.DB.QueryRow(countQuery, taskID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := `SELECT id, user_id, task_id, action, entity_type, entity_id, details, ip_address, user_agent, created_at 
	          FROM activity_logs 
	          WHERE task_id = $1 
	          ORDER BY created_at DESC 
	          LIMIT $2 OFFSET $3`

	rows, err := m.DB.Query(query, taskID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	logs := make([]ActivityLog, 0)
	for rows.Next() {
		var log ActivityLog
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.TaskID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.Details,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (m *ActivityLogModel) GetByID(id int64) (*ActivityLog, error) {
	query := `SELECT id, user_id, task_id, action, entity_type, entity_id, details, ip_address, user_agent, created_at 
	          FROM activity_logs WHERE id = $1`
	var log ActivityLog
	err := m.DB.QueryRow(query, id).Scan(
		&log.ID,
		&log.UserID,
		&log.TaskID,
		&log.Action,
		&log.EntityType,
		&log.EntityID,
		&log.Details,
		&log.IPAddress,
		&log.UserAgent,
		&log.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &log, nil
}
