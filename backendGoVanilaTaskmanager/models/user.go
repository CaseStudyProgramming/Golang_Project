package models

import (
	"database/sql"
	"taskmanager/utils"
)

// User entity
type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Timezone     string `json:"timezone"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

// UserModelInterface defines the interface for user model operations
type UserModelInterface interface {
	Create(user *User) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id int64) (*User, error)
	Update(user *User) error
	UpdateTimezone(userID int64, timezone string) error
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

func (m *UserModel) Create(user *User) (*User, error) {
	query := `INSERT INTO users (name, email, password_hash, timezone) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, name, email, timezone, created_at, updated_at`
	err := m.DB.QueryRow(query, user.Name, user.Email, user.PasswordHash, user.Timezone).
		Scan(&user.ID, &user.Name, &user.Email, &user.Timezone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, name, email, password_hash, timezone, created_at, updated_at 
	          FROM users WHERE email = $1`
	var user User
	err := m.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetByID(id int64) (*User, error) {
	query := `SELECT id, name, email, password_hash, timezone, created_at, updated_at 
	          FROM users WHERE id = $1`
	var user User
	err := m.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE users 
	          SET name = $1, email = $2, timezone = $3, updated_at = $4 
	          WHERE id = $5 
	          RETURNING updated_at`
	err := m.DB.QueryRow(query, user.Name, user.Email, user.Timezone, currentTime, user.ID).Scan(&user.UpdatedAt)
	return err
}

func (m *UserModel) UpdateTimezone(userID int64, timezone string) error {
	currentTime := utils.CurrentEpochMillis()
	query := `UPDATE users 
	          SET timezone = $1, updated_at = $2 
	          WHERE id = $3`
	_, err := m.DB.Exec(query, timezone, currentTime, userID)
	return err
}
