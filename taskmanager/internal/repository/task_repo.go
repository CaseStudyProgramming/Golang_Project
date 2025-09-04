package repository

import (
	"database/sql"
	"taskmanager/internal/entity"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) Create(task *entity.Task) error {
	query := `INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING id, created_at`
	return r.DB.QueryRow(query, task.Title, task.Completed).Scan(&task.ID, &task.CreatedAt)
}

func (r *TaskRepository) GetAll() ([]entity.Task, error) {
	query := `SELECT id, title, completed, created_at FROM tasks ORDER BY id ASC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
