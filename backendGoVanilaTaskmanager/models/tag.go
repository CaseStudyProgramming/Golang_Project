package models

import (
	"database/sql"
)

// Tag entity
type Tag struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	ColorHex  string `json:"color_hex"`
	CreatedAt int64  `json:"created_at"`
}

// TaskTag represents the junction table for many-to-many relationship
type TaskTag struct {
	TaskID    int64 `json:"task_id"`
	TagID     int64 `json:"tag_id"`
	CreatedAt int64 `json:"created_at"`
}

type TagModel struct {
	DB *sql.DB
}

// TagModelInterface defines the interface for tag model operations
type TagModelInterface interface {
	Create(tag *Tag) (*Tag, error)
	GetAll(userID int64) ([]Tag, error)
	GetByID(userID int64, id int64) (*Tag, error)
	Update(userID int64, tag *Tag) error
	Delete(userID int64, id int64) error
	// Task-tag relationship operations
	AddTagToTask(taskID int64, tagID int64) error
	RemoveTagFromTask(taskID int64, tagID int64) error
	GetTagsByTaskID(taskID int64) ([]Tag, error)
	GetTasksByTagID(tagID int64) ([]Task, error)
}

func NewTagModel(db *sql.DB) *TagModel {
	return &TagModel{DB: db}
}

func (m *TagModel) Create(tag *Tag) (*Tag, error) {
	query := `INSERT INTO tags (user_id, name, color_hex) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, user_id, name, color_hex, created_at`
	err := m.DB.QueryRow(query, tag.UserID, tag.Name, tag.ColorHex).
		Scan(&tag.ID, &tag.UserID, &tag.Name, &tag.ColorHex, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (m *TagModel) GetAll(userID int64) ([]Tag, error) {
	query := `SELECT id, user_id, name, color_hex, created_at 
	          FROM tags 
	          WHERE user_id = $1 
	          ORDER BY created_at DESC`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m *TagModel) GetByID(userID int64, id int64) (*Tag, error) {
	query := `SELECT id, user_id, name, color_hex, created_at 
	          FROM tags WHERE id = $1 AND user_id = $2`
	var tag Tag
	err := m.DB.QueryRow(query, id, userID).Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Name,
		&tag.ColorHex,
		&tag.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (m *TagModel) Update(userID int64, tag *Tag) error {
	query := `UPDATE tags 
	          SET name = $1, color_hex = $2 
	          WHERE id = $3 AND user_id = $4`
	_, err := m.DB.Exec(query, tag.Name, tag.ColorHex, tag.ID, userID)
	return err
}

func (m *TagModel) Delete(userID int64, id int64) error {
	query := `DELETE FROM tags WHERE id = $1 AND user_id = $2`
	_, err := m.DB.Exec(query, id, userID)
	return err
}

// Task-tag relationship operations
func (m *TagModel) AddTagToTask(taskID int64, tagID int64) error {
	query := `INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(query, taskID, tagID)
	return err
}

func (m *TagModel) RemoveTagFromTask(taskID int64, tagID int64) error {
	query := `DELETE FROM task_tags WHERE task_id = $1 AND tag_id = $2`
	_, err := m.DB.Exec(query, taskID, tagID)
	return err
}

func (m *TagModel) GetTagsByTaskID(taskID int64) ([]Tag, error) {
	query := `SELECT t.id, t.user_id, t.name, t.color_hex, t.created_at 
	          FROM tags t 
	          INNER JOIN task_tags tt ON t.id = tt.tag_id 
	          WHERE tt.task_id = $1 
	          ORDER BY t.created_at DESC`

	rows, err := m.DB.Query(query, taskID)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m *TagModel) GetTasksByTagID(tagID int64) ([]Task, error) {
	query := `SELECT t.id, t.user_id, t.category_id, t.title, t.sub_title, t.description, t.completed, t.due_date, t.priority, t.created_at, t.updated_at, t.deleted_at 
	          FROM tasks t 
	          INNER JOIN task_tags tt ON t.id = tt.task_id 
	          WHERE tt.tag_id = $1 AND t.deleted_at IS NULL 
	          ORDER BY t.created_at DESC`

	rows, err := m.DB.Query(query, tagID)
	if err != nil {
		return nil, err
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
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
