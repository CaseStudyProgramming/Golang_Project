package models

import (
	"database/sql"
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
	query := `INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING id, title, completed, created_at`
	err := m.DB.QueryRow(query, task.Title, task.Completed).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *TaskModel) GetAll(completed *bool, offset int, limit int, search string) ([]Task, int, error) {
	// Get total count
	var countQuery string
	var countArgs []interface{}
	var total int

	if completed == nil {
		countQuery = `SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL`
	} else {
		countQuery = `SELECT COUNT(*) FROM tasks WHERE completed = $1 AND deleted_at IS NULL`
		countArgs = append(countArgs, *completed)
	}

	if search != "" {
		if completed == nil {
			countQuery += " AND (title LIKE $1 OR description LIKE $1)"
			countArgs = append(countArgs, "%"+search+"%")
		} else {
			countQuery += " AND (title LIKE $2 OR description LIKE $2)"
			countArgs = append(countArgs, "%"+search+"%")
		}
	}

	if err := m.DB.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated data
	var (
		query string
		args  []interface{}
	)

	if completed == nil {
		query = `
			SELECT id, title, completed, created_at, deleted_at
			FROM tasks
			WHERE deleted_at IS NULL AND (title LIKE $1 OR description LIKE $1)
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = append(args, "%"+search+"%", limit, offset)
	} else {
		query = `
			SELECT id, title, completed, created_at, deleted_at
			FROM tasks
			WHERE completed = $1 AND deleted_at IS NULL AND (title LIKE $2 OR description LIKE $2)
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		args = append(args, *completed, "%"+search+"%", limit, offset)
	}

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
			&task.Completed,
			&task.CreatedAt,
			&task.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func (m *TaskModel) GetByID(id int64) (*Task, error) {
	query := `SELECT id, title, completed, created_at FROM tasks WHERE id = $1 AND deleted_at IS NULL`
	var task Task
	if err := m.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
		return nil, err
	}
	return &task, nil
}

func (m *TaskModel) Update(task *Task) error {
	query := `UPDATE tasks SET title = $1, completed = $2 WHERE id = $3`
	_, err := m.DB.Exec(query, task.Title, task.Completed, task.ID)
	return err
}

func (m *TaskModel) Complete(id int64) error {
	query := `UPDATE tasks SET completed = true WHERE id = $1`
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
	query := `UPDATE tasks SET deleted_at = NULL WHERE id = $1 AND deleted_at IS NOT NULL`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *TaskModel) MarkTaskAsCompleted(id int64) error {
	query := `UPDATE tasks SET completed = true WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *TaskModel) MarkTaskAsUncompleted(id int64) error {
	query := `UPDATE tasks SET completed = false WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}
