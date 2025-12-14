package postgres

import "time"

type User struct {
	ID           string    `db:"id" json:"id"` // uuid string
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	FullName     string    `db:"full_name" json:"full_name"`
	RoleID       string    `db:"role_id" json:"role_id"` // FK -> roles.id
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
