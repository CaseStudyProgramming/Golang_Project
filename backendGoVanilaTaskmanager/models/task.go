package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Task entity
type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	SubTitle    string     `json:"sub_title,omitempty"`
	Description string     `json:"description,omitempty"`
	Completed   bool       `json:"completed"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type TaskModel struct {
	DB *sql.DB
}

func NewTaskModel(db *sql.DB) *TaskModel {
	return &TaskModel{DB: db}
}

func (m *TaskModel) Create(task *Task) (*Task, error) {
	query := `INSERT INTO tasks (title, sub_title, description, completed, due_date) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, title, sub_title, description, completed, due_date, created_at, updated_at`
	err := m.DB.QueryRow(query, task.Title, task.SubTitle, task.Description, task.Completed, task.DueDate).
		Scan(&task.ID, &task.Title, &task.SubTitle, &task.Description, &task.Completed, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *TaskModel) GetAll(completed *bool, offset int, limit int, search string) ([]Task, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL`
	var countArgs []interface{}

	if completed != nil {
		countQuery += ` AND completed = $1`
		countArgs = append(countArgs, *completed)
	}

	if search != "" {
		placeholder := fmt.Sprintf("$%d", len(countArgs)+1)
		countQuery += fmt.Sprintf(` AND (title LIKE %s OR description LIKE %s)`, placeholder, placeholder)
		countArgs = append(countArgs, "%"+search+"%")
	}

	var total int
	if err := m.DB.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := `SELECT id, title, sub_title, description, completed, due_date, created_at, updated_at, deleted_at 
	          FROM tasks 
	          WHERE deleted_at IS NULL`
	var args []interface{}

	if completed != nil {
		query += ` AND completed = $1`
		args = append(args, *completed)
	}

	if search != "" {
		placeholder := fmt.Sprintf("$%d", len(args)+1)
		query += fmt.Sprintf(` AND (title LIKE %s OR description LIKE %s)`, placeholder, placeholder)
		args = append(args, "%"+search+"%")
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.SubTitle,
			&task.Description,
			&task.Completed,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func (m *TaskModel) GetByID(id int64) (*Task, error) {
	query := `SELECT id, title, sub_title, description, completed, due_date, created_at, updated_at, deleted_at 
	          FROM tasks WHERE id = $1 AND deleted_at IS NULL`
	var task Task
	err := m.DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.SubTitle,
		&task.Description,
		&task.Completed,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (m *TaskModel) Update(task *Task) error {
	query := `UPDATE tasks 
	          SET title = $1, sub_title = $2, description = $3, completed = $4, due_date = $5, updated_at = NOW() 
	          WHERE id = $6 
	          RETURNING updated_at`
	err := m.DB.QueryRow(query, task.Title, task.SubTitle, task.Description, task.Completed, task.DueDate, task.ID).Scan(&task.UpdatedAt)
	return err
}

func (m *TaskModel) Complete(id int64) error {
	query := `UPDATE tasks SET completed = true, updated_at = NOW() WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *TaskModel) Delete(id int64, softDelete bool) error {
	if softDelete {
		query := `UPDATE tasks SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
		_, err := m.DB.Exec(query, id)
		return err
	} else {
		query := `DELETE FROM tasks WHERE id = $1 AND deleted_at IS NULL`
		_, err := m.DB.Exec(query, id)
		return err
	}
}

func (m *TaskModel) RestoreTask(id int64) error {
	query := `UPDATE tasks SET deleted_at = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *TaskModel) MarkTaskAsCompleted(id int64) error {
	query := `UPDATE tasks SET completed = true, updated_at = NOW() WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *TaskModel) MarkTaskAsUncompleted(id int64) error {
	query := `UPDATE tasks SET completed = false, updated_at = NOW() WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
