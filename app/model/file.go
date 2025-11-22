package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID       string             `json:"user_id" bson:"user_id"`
	FileName     string             `json:"file_name" bson:"file_name"`
	OriginalName string             `json:"original_name" bson:"original_name"`
	FilePath     string             `json:"file_path" bson:"file_path"`
	FileSize     int64              `json:"file_size" bson:"file_size"`
	FileType     string             `json:"file_type" bson:"file_type"`
	Category     string             `json:"category" bson:"category"` // "photo" atau "certificate"
	UploadedAt   time.Time          `json:"uploaded_at" bson:"uploaded_at"`
	UploadedBy   string             `json:"uploaded_by" bson:"uploaded_by"` // Admin atau User ID
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt    *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type FileResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	Category     string    `json:"category"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UploadedBy   UserInfo  `json:"uploaded_by"` // Contains username, email, role
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UploadPhotoRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type UploadCertificateRequest struct {
	UserID string `json:"user_id" validate:"required"`
}
