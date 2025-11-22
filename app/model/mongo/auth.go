package model

import (

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
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
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AlumniJWTClaims struct {
	AlumniID string `json:"alumni_id"`
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
