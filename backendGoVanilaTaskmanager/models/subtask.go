package models

import (
	"database/sql"
	"time"
)

// Subtask entity
type Subtask struct {
	ID          int64     `json:"id"`
	TaskID      int64     `json:"task_id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SubtaskModel struct {
	DB *sql.DB
}

// SubtaskModelInterface defines the interface for subtask model operations
type SubtaskModelInterface interface {
	Create(subtask *Subtask) (*Subtask, error)
	GetByTaskID(taskID int64) ([]Subtask, error)
	GetByID(id int64) (*Subtask, error)
	Update(subtask *Subtask) error
	Delete(id int64) error
	Toggle(id int64) error
	CalculateProgress(taskID int64) (int, error)
}

func NewSubtaskModel(db *sql.DB) *SubtaskModel {
	return &SubtaskModel{DB: db}
}

func (m *SubtaskModel) Create(subtask *Subtask) (*Subtask, error) {
	query := `INSERT INTO subtasks (task_id, title, is_completed) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, task_id, title, is_completed, created_at, updated_at`
	err := m.DB.QueryRow(query, subtask.TaskID, subtask.Title, subtask.IsCompleted).
		Scan(&subtask.ID, &subtask.TaskID, &subtask.Title, &subtask.IsCompleted, &subtask.CreatedAt, &subtask.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return subtask, nil
}

func (m *SubtaskModel) GetByTaskID(taskID int64) ([]Subtask, error) {
	query := `SELECT id, task_id, title, is_completed, created_at, updated_at 
	          FROM subtasks 
	          WHERE task_id = $1 
	          ORDER BY created_at ASC`

	rows, err := m.DB.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subtasks := make([]Subtask, 0)
	for rows.Next() {
		var subtask Subtask
		if err := rows.Scan(
			&subtask.ID,
			&subtask.TaskID,
			&subtask.Title,
			&subtask.IsCompleted,
			&subtask.CreatedAt,
			&subtask.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subtasks = append(subtasks, subtask)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subtasks, nil
}

func (m *SubtaskModel) GetByID(id int64) (*Subtask, error) {
	query := `SELECT id, task_id, title, is_completed, created_at, updated_at 
	          FROM subtasks WHERE id = $1`
	var subtask Subtask
	err := m.DB.QueryRow(query, id).Scan(
		&subtask.ID,
		&subtask.TaskID,
		&subtask.Title,
		&subtask.IsCompleted,
		&subtask.CreatedAt,
		&subtask.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &subtask, nil
}

func (m *SubtaskModel) Update(subtask *Subtask) error {
	query := `UPDATE subtasks 
	          SET title = $1, is_completed = $2, updated_at = NOW() 
	          WHERE id = $3
	          RETURNING updated_at`
	err := m.DB.QueryRow(query, subtask.Title, subtask.IsCompleted, subtask.ID).Scan(&subtask.UpdatedAt)
	return err
}

func (m *SubtaskModel) Delete(id int64) error {
	query := `DELETE FROM subtasks WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *SubtaskModel) Toggle(id int64) error {
	query := `UPDATE subtasks 
	          SET is_completed = NOT is_completed, updated_at = NOW() 
	          WHERE id = $1
	          RETURNING is_completed, updated_at`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *SubtaskModel) CalculateProgress(taskID int64) (int, error) {
	query := `SELECT 
	          COUNT(*) as total,
	          SUM(CASE WHEN is_completed = true THEN 1 ELSE 0 END) as completed
	          FROM subtasks 
	          WHERE task_id = $1`

	var total, completed int
	err := m.DB.QueryRow(query, taskID).Scan(&total, &completed)
	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 0, nil
	}

	progress := (completed * 100) / total
	return progress, nil
}
