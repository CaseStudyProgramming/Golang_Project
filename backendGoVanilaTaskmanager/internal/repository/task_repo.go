package repository

import (
	"database/sql"
	"taskmanager/internal/entity"
)

type TaskRepository struct {
	DB *sql.DB
}

func (r *TaskRepository) Uncomplete(id int64) any {
	panic("unimplemented")
}

// CREATE REPOSITORY
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// POST REPOSITORY

func (r *TaskRepository) Create(task *entity.Task) (*entity.Task, error) {
	query := `INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING id, title, completed, created_at`
	err := r.DB.QueryRow(query, task.Title, task.Completed).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// GET ALL DATA REPOSITORY

func (r *TaskRepository) GetAll(completed *bool, offset int, limit int, search string) ([]entity.Task, int, error) {
	// Get total count
	var countQuery string
	var countArgs []interface{}
	var total int

	if completed == nil {
		countQuery = `SELECT COUNT(*) FROM tasks`
	} else {
		countQuery = `SELECT COUNT(*) FROM tasks WHERE completed = $1`
		countArgs = append(countArgs, *completed)
	}

	if search != "" {
		if completed == nil {
			countQuery += " WHERE (title LIKE $1 OR description LIKE $1)"
			countArgs = append(countArgs, "%"+search+"%")
		} else {
			countQuery += " AND (title LIKE $2 OR description LIKE $2)"
			countArgs = append(countArgs, "%"+search+"%")
		}
	}

	if err := r.DB.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated data
	var (
		query string
		args  []interface{}
	)

	if completed == nil {
		query = `
			SELECT id, title, completed, created_at
			FROM tasks
			WHERE (title LIKE $1 OR description LIKE $1)
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = append(args, "%"+search+"%", limit, offset)
	} else {
		query = `
			SELECT id, title, completed, created_at
			FROM tasks
			WHERE completed = $1 AND (title LIKE $2 OR description LIKE $2)
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		args = append(args, *completed, "%"+search+"%", limit, offset)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]entity.Task, 0)
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
			&task.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
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
	query := `UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP() WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err

// MarkTaskAsCompleted
func (r *TaskRepository) MarkTaskAsCompleted(id int64) error {
	query := `UPDATE tasks SET completed = true WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

// MarkTaskAsUncompleted
func (r *TaskRepository) MarkTaskAsUncompleted(id int64) error {
	query := `UPDATE tasks SET completed = false WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
