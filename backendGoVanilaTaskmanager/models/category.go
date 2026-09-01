package models

import (
	"database/sql"
	"time"
)

// Category entity
type Category struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	ColorHex  string    `json:"color_hex"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryModel struct {
	DB *sql.DB
}

// CategoryModelInterface defines the interface for category model operations
type CategoryModelInterface interface {
	Create(category *Category) (*Category, error)
	GetAll(userID int64) ([]Category, error)
	GetByID(userID int64, id int64) (*Category, error)
	Update(userID int64, category *Category) error
	Delete(userID int64, id int64) error
}

func NewCategoryModel(db *sql.DB) *CategoryModel {
	return &CategoryModel{DB: db}
}

func (m *CategoryModel) Create(category *Category) (*Category, error) {
	query := `INSERT INTO categories (user_id, name, color_hex) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, user_id, name, color_hex, created_at`
	err := m.DB.QueryRow(query, category.UserID, category.Name, category.ColorHex).
		Scan(&category.ID, &category.UserID, &category.Name, &category.ColorHex, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (m *CategoryModel) GetAll(userID int64) ([]Category, error) {
	query := `SELECT id, user_id, name, color_hex, created_at 
	          FROM categories 
	          WHERE user_id = $1 
	          ORDER BY created_at DESC`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		if err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.ColorHex,
			&category.CreatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (m *CategoryModel) GetByID(userID int64, id int64) (*Category, error) {
	query := `SELECT id, user_id, name, color_hex, created_at 
	          FROM categories WHERE id = $1 AND user_id = $2`
	var category Category
	err := m.DB.QueryRow(query, id, userID).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.ColorHex,
		&category.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (m *CategoryModel) Update(userID int64, category *Category) error {
	query := `UPDATE categories 
	          SET name = $1, color_hex = $2 
	          WHERE id = $3 AND user_id = $4`
	_, err := m.DB.Exec(query, category.Name, category.ColorHex, category.ID, userID)
	return err
}

func (m *CategoryModel) Delete(userID int64, id int64) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	_, err := m.DB.Exec(query, id, userID)
	return err
}
