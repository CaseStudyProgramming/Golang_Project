package models

import (
	"database/sql"
	"fmt"
	"taskmanager/utils"
)

// Priority represents task priority levels
type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)

// Task entity
type Task struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	CategoryID         *int64    `json:"category_id,omitempty"`
	Title              string    `json:"title"`
	SubTitle           string    `json:"sub_title,omitempty"`
	Description        string    `json:"description,omitempty"`
	Completed          bool      `json:"completed"`
	DueDate            *int64    `json:"due_date,omitempty"`
	Priority           Priority  `json:"priority"`
	CreatedAt          int64     `json:"created_at"`
	UpdatedAt          int64     `json:"updated_at"`
	DeletedAt          *int64    `json:"deleted_at"`
	Tags               []Tag     `json:"tags,omitempty"`
	ProgressPercentage int       `json:"progress_percentage"`
	Subtasks           []Subtask `json:"subtasks,omitempty"`
}

type TaskModel struct {
	DB *sql.DB
}

// TaskModelInterface defines the interface for task model operations
type TaskModelInterface interface {
	Create(task *Task) (*Task, error)
	GetAll(userID int64, completed *bool, offset, limit int, search string, priority *Priority, categoryID *int64, sortBy string, sortOrder string) ([]Task, int, error)
	GetByID(userID int64, id int64) (*Task, error)
	Update(userID int64, task *Task) error
	Complete(userID int64, id int64) error
	Delete(userID int64, id int64, softDelete bool) error
	RestoreTask(userID int64, id int64) error
	MarkTaskAsCompleted(userID int64, id int64) error
	MarkTaskAsUncompleted(userID int64, id int64) error
	LoadTags(task *Task) error
	AddTagToTask(taskID int64, tagID int64) error
	RemoveTagFromTask(taskID int64, tagID int64) error
	LoadSubtasks(task *Task) error
	UpdateProgress(taskID int64) error
	GetAnalyticsSummary(userID int64) (map[string]interface{}, error)
	BulkDelete(userID int64, taskIDs []int64) error
	BulkComplete(userID int64, taskIDs []int64) error
}

func NewTaskModel(db *sql.DB) *TaskModel {
	return &TaskModel{DB: db}
}

func (m *TaskModel) Create(task *Task) (*Task, error) {
	query := `INSERT INTO tasks (user_id, category_id, title, sub_title, description, completed, due_date, priority, progress_percentage) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0) 
	          RETURNING id, user_id, category_id, title, sub_title, description, completed, due_date, priority, created_at, updated_at, progress_percentage`
	err := m.DB.QueryRow(query, task.UserID, task.CategoryID, task.Title, task.SubTitle, task.Description, task.Completed, task.DueDate, task.Priority).
		Scan(&task.ID, &task.UserID, &task.CategoryID, &task.Title, &task.SubTitle, &task.Description, &task.Completed, &task.DueDate, &task.Priority, &task.CreatedAt, &task.UpdatedAt, &task.ProgressPercentage)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *TaskModel) GetAll(userID int64, completed *bool, offset int, limit int, search string, priority *Priority, categoryID *int64, sortBy string, sortOrder string) ([]Task, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL AND user_id = $1`
	var countArgs []interface{}
	countArgs = append(countArgs, userID)

	if completed != nil {
		countQuery += ` AND completed = $2`
		countArgs = append(countArgs, *completed)
	}

	if priority != nil {
		countQuery += ` AND priority = $` + fmt.Sprintf("%d", len(countArgs)+1)
		countArgs = append(countArgs, *priority)
	}

	if categoryID != nil {
		countQuery += ` AND category_id = $` + fmt.Sprintf("%d", len(countArgs)+1)
		countArgs = append(countArgs, *categoryID)
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
	query := `SELECT id, user_id, category_id, title, sub_title, description, completed, due_date, priority, created_at, updated_at, deleted_at, progress_percentage 
	          FROM tasks 
	          WHERE deleted_at IS NULL AND user_id = $1`
	var args []interface{}
	args = append(args, userID)

	if completed != nil {
		query += ` AND completed = $2`
		args = append(args, *completed)
	}

	if priority != nil {
		query += ` AND priority = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *priority)
	}

	if categoryID != nil {
		query += ` AND category_id = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *categoryID)
	}

	if search != "" {
		placeholder := fmt.Sprintf("$%d", len(args)+1)
		query += fmt.Sprintf(` AND (title LIKE %s OR description LIKE %s)`, placeholder, placeholder)
		args = append(args, "%"+search+"%")
	}

	// Handle sorting
	orderBy := "created_at DESC"
	if sortBy != "" {
		validSortFields := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"due_date":   true,
			"priority":   true,
			"title":      true,
		}
		if validSortFields[sortBy] {
			order := "ASC"
			if sortOrder == "DESC" {
				order = "DESC"
			}
			orderBy = sortBy + " " + order
		}
	}

	query += fmt.Sprintf(` ORDER BY %s LIMIT $%d OFFSET $%d`, orderBy, len(args)+1, len(args)+2)
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
			&task.UserID,
			&task.CategoryID,
			&task.Title,
			&task.SubTitle,
			&task.Description,
			&task.Completed,
			&task.DueDate,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DeletedAt,
			&task.ProgressPercentage,
		); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (m *TaskModel) GetByID(userID int64, id int64) (*Task, error) {
	query := `SELECT id, user_id, category_id, title, sub_title, description, completed, due_date, priority, created_at, updated_at, deleted_at, progress_percentage 
	          FROM tasks WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	var task Task
	err := m.DB.QueryRow(query, id, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.CategoryID,
		&task.Title,
		&task.SubTitle,
		&task.Description,
		&task.Completed,
		&task.DueDate,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
		&task.ProgressPercentage,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (m *TaskModel) Update(userID int64, task *Task) error {
	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE tasks 
	          SET category_id = $1, title = $2, sub_title = $3, description = $4, completed = $5, due_date = $6, priority = $7, updated_at = $8 
	          WHERE id = $9 AND user_id = $10
	          RETURNING updated_at`
	err := m.DB.QueryRow(query, task.CategoryID, task.Title, task.SubTitle, task.Description, task.Completed, task.DueDate, task.Priority, currentTime, task.ID, userID).Scan(&task.UpdatedAt)
	return err
}

func (m *TaskModel) Complete(userID int64, id int64) error {
	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE tasks SET completed = true, updated_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := m.DB.Exec(query, currentTime, id, userID)
	return err
}

func (m *TaskModel) Delete(userID int64, id int64, softDelete bool) error {
	if softDelete {
		currentTime := utils.CurrentEpochMillis()
		query := `UPDATE tasks SET deleted_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL`
		_, err := m.DB.Exec(query, currentTime, id, userID)
		return err
	} else {
		query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
		_, err := m.DB.Exec(query, id, userID)
		return err
	}
}

func (m *TaskModel) RestoreTask(userID int64, id int64) error {
	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE tasks SET deleted_at = NULL, updated_at = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NOT NULL`
	_, err := m.DB.Exec(query, currentTime, id, userID)
	return err
}

func (m *TaskModel) MarkTaskAsCompleted(userID int64, id int64) error {
	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE tasks SET completed = true, updated_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := m.DB.Exec(query, currentTime, id, userID)
	return err
}

func (m *TaskModel) MarkTaskAsUncompleted(userID int64, id int64) error {
	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE tasks SET completed = false, updated_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := m.DB.Exec(query, currentTime, id, userID)
	return err
}

func (m *TaskModel) LoadTags(task *Task) error {
	query := `SELECT t.id, t.user_id, t.name, t.color_hex, t.created_at 
	          FROM tags t 
	          INNER JOIN task_tags tt ON t.id = tt.tag_id 
	          WHERE tt.task_id = $1 
	          ORDER BY t.created_at DESC`

	rows, err := m.DB.Query(query, task.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	tags := make([]Tag, 0)
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(
			&tag.ID,
			&tag.UserID,
			&tag.Name,
			&tag.ColorHex,
			&tag.CreatedAt,
		); err != nil {
			return err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	task.Tags = tags
	return nil
}

func (m *TaskModel) AddTagToTask(taskID int64, tagID int64) error {
	query := `INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(query, taskID, tagID)
	return err
}

func (m *TaskModel) RemoveTagFromTask(taskID int64, tagID int64) error {
	query := `DELETE FROM task_tags WHERE task_id = $1 AND tag_id = $2`
	_, err := m.DB.Exec(query, taskID, tagID)
	return err
}

func (m *TaskModel) LoadSubtasks(task *Task) error {
	query := `SELECT id, task_id, title, is_completed, created_at, updated_at 
	          FROM subtasks 
	          WHERE task_id = $1 
	          ORDER BY created_at ASC`

	rows, err := m.DB.Query(query, task.ID)
	if err != nil {
		return err
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
			return err
		}
		subtasks = append(subtasks, subtask)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	task.Subtasks = subtasks
	return nil
}

func (m *TaskModel) UpdateProgress(taskID int64) error {
	subtaskModel := NewSubtaskModel(m.DB)
	progress, err := subtaskModel.CalculateProgress(taskID)
	if err != nil {
		return err
	}

	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE tasks SET progress_percentage = $1, updated_at = $2 WHERE id = $3`
	_, err = m.DB.Exec(query, progress, currentTime, taskID)
	return err
}

func (m *TaskModel) GetAnalyticsSummary(userID int64) (map[string]interface{}, error) {
	summary := make(map[string]interface{})
	now := utils.CurrentEpochMillis()

	// Total active tasks
	var totalActive int
	err := m.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND deleted_at IS NULL`, userID).Scan(&totalActive)
	if err != nil {
		return nil, err
	}
	summary["total_active"] = totalActive

	// Total completed tasks
	var totalCompleted int
	err = m.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND deleted_at IS NULL AND completed = true`, userID).Scan(&totalCompleted)
	if err != nil {
		return nil, err
	}
	summary["total_completed"] = totalCompleted

	// Total overdue tasks (not completed and due_date is in the past)
	var totalOverdue int
	err = m.DB.QueryRow(`SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND deleted_at IS NULL AND completed = false AND due_date < $2`, userID, now).Scan(&totalOverdue)
	if err != nil {
		return nil, err
	}
	summary["total_overdue"] = totalOverdue

	// Completion percentage
	completionPercentage := 0.0
	if totalActive > 0 {
		completionPercentage = float64(totalCompleted) / float64(totalActive) * 100
	}
	summary["completion_percentage"] = completionPercentage

	// Priority distribution
	priorityQuery := `SELECT priority, COUNT(*) FROM tasks WHERE user_id = $1 AND deleted_at IS NULL GROUP BY priority`
	rows, err := m.DB.Query(priorityQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	priorityDistribution := make(map[string]int)
	for rows.Next() {
		var priority string
		var count int
		if err := rows.Scan(&priority, &count); err != nil {
			return nil, err
		}
		priorityDistribution[priority] = count
	}
	summary["priority_distribution"] = priorityDistribution

	return summary, nil
}

func (m *TaskModel) BulkDelete(userID int64, taskIDs []int64) error {
	if len(taskIDs) == 0 {
		return nil
	}

	currentTime := utils.CurrentEpochMillis()
	// Build the query with IN clause
	query := `UPDATE tasks SET deleted_at = $1 WHERE user_id = $2 AND id = ANY($3) AND deleted_at IS NULL`
	_, err := m.DB.Exec(query, currentTime, userID, taskIDs)
	return err
}

func (m *TaskModel) BulkComplete(userID int64, taskIDs []int64) error {
	if len(taskIDs) == 0 {
		return nil
	}

	currentTime := utils.CurrentEpochMillis()
	// Build the query with IN clause
	query := `UPDATE tasks SET completed = true, updated_at = $1 WHERE user_id = $2 AND id = ANY($3) AND deleted_at IS NULL`
	_, err := m.DB.Exec(query, currentTime, userID, taskIDs)
	return err
}
