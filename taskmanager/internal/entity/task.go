package entity

import "time"

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
}
