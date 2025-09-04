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
