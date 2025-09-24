package repository

import (
	"database/sql"
	"taskmanager/internal/entity"
)

type TaskRepository struct {
	DB *sql.DB
}

// CREATE REPOSITORY
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// POST REPO

func (r *TaskRepository) Create(task *entity.Task) (*entity.Task, error) {
	query := `INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING *`
	return r.DB.QueryRow(query, task.Title, task.Completed).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)

// GET ALL DATA

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

// GET BY ID
func (r *TaskRepository) GetByID(id int64) (*entity.Task, error) {
	query := `SELECT id, title, completed, created_at FROM tasks WHERE id = $1`
	var task entity.Task
	if err := r.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
		return nil, err
	}
	return &task, nil
}

// PUT
func (r *TaskRepository) Update(task *entity.Task) error {
	query := `UPDATE tasks SET title = $1, completed = $2 WHERE id = $3`
	_, err := r.DB.Exec(query, task.Title, task.Completed, task.ID)
	return err
}

// PATCH
func (r *TaskRepository) Complete(id int64) error {
	query := `UPDATE tasks SET completed = true WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

// DELETE
func (r *TaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
