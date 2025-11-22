package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	AlumniID int    `json:"alumni_id"`
	Username string `json:"username"`
	NIM      string `json:"nim"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AlumniJWTClaims struct {
	AlumniID int    `json:"alumni_id"`
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
