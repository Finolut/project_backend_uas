package service

import (
	"clean-arch/app/model/mongo"
	"clean-arch/app/repository/mongo"
	"clean-arch/utils/mongo"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// File size limits
	maxPhotoSize       = 1 * 1024 * 1024 // 1MB
	maxCertificateSize = 2 * 1024 * 1024 // 2MB
	uploadBasePath     = "./uploads"
	photosDir          = "photos"
	certificatesDir    = "certificates"
)

// UploadPhotoService handles photo upload
func UploadPhotoService(c *fiber.Ctx, db *mongo.Database) error {
	currentUserID := c.Locals("user_id")
	currentRole := c.Locals("role")

	if currentUserID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID not found in token",
		})
	}

	userID := currentUserID.(string)

	targetUserID := c.FormValue("user_id")
	if targetUserID != "" {
		// Only admin can upload for another user
		if currentRole != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Only admin can upload for other users",
			})
		}
		userID = targetUserID
	}

	// Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File is required",
			"error":   err.Error(),
		})
	}

	if fileHeader.Size > maxPhotoSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Photo size must not exceed 1MB",
		})
	}

	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedMimeTypes[contentType] {
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Only JPEG and PNG formats are allowed",
			})
		}
	}

	_, err = repository.GetUserByID(db, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	uploadedFile, err := saveFile(db, fileHeader, "photo", userID, currentUserID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to upload file",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Photo uploaded successfully",
		"data":    toFileResponse(uploadedFile, db),
	})
}

// UploadCertificateService handles certificate/diploma upload
func UploadCertificateService(c *fiber.Ctx, db *mongo.Database) error {
	currentUserID := c.Locals("user_id")
	currentRole := c.Locals("role")

	if currentUserID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User ID not found in token",
		})
	}

	userID := currentUserID.(string)

	targetUserID := c.FormValue("user_id")
	if targetUserID != "" {
		// Only admin can upload for another user
		if currentRole != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Only admin can upload for other users",
			})
		}
		userID = targetUserID
	}

	// Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File is required",
			"error":   err.Error(),
		})
	}

	if fileHeader.Size > maxCertificateSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Certificate size must not exceed 2MB",
		})
	}

	contentType := fileHeader.Header.Get("Content-Type")
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	if contentType != "application/pdf" && ext != ".pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Only PDF format is allowed",
		})
	}

	_, err = repository.GetUserByID(db, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	uploadedFile, err := saveFile(db, fileHeader, "certificate", userID, currentUserID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to upload file",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Certificate uploaded successfully",
		"data":    toFileResponse(uploadedFile, db),
	})
}

// GetFilesService retrieves files for specific user
func GetFilesService(c *fiber.Ctx, db *mongo.Database) error {
	userID := c.Query("user_id")
	category := c.Query("category") // "photo" atau "certificate"

	if userID == "" || category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "user_id and category are required",
		})
	}

	files, err := repository.GetFileByUserID(db, userID, category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve files",
			"error":   err.Error(),
		})
	}

	var responses []model.FileResponse
	for _, file := range files {
		responses = append(responses, *toFileResponse(&file, db))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Files retrieved successfully",
		"data":    responses,
	})
}

// DeleteFileService soft deletes a file
func DeleteFileService(c *fiber.Ctx, db *mongo.Database) error {
	fileID := c.Params("id")

	file, err := repository.GetFileByID(db, fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
		})
	}

	currentUserID := c.Locals("user_id")
	currentRole := c.Locals("role")

	if currentRole != "admin" && currentUserID != file.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "You can only delete your own files",
		})
	}

	// Soft delete
	err = repository.DeleteFile(db, fileID, currentUserID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete file",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

// Helper function to save file to disk and database
func saveFile(db *mongo.Database, fileHeader *multipart.FileHeader, category, userID, uploadedBy string) (*model.File, error) {
	uploadDir := filepath.Join(uploadBasePath, category)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFileName)

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		return nil, err
	}

	fileModel := &model.File{
		UserID:       userID,
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     fileHeader.Header.Get("Content-Type"),
		Category:     category,
		UploadedAt:   utils.GetNowTime(),
		UploadedBy:   uploadedBy,
	}

	if err := repository.CreateFile(db, fileModel); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return fileModel, nil
}

// toFileResponse converts File model to FileResponse
func toFileResponse(file *model.File, db *mongo.Database) *model.FileResponse {
	// Fetch user info from users collection
	userInfo := getUserInfo(db, file.UploadedBy)

	return &model.FileResponse{
		ID:           file.ID.Hex(),
		UserID:       file.UserID,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		Category:     file.Category,
		UploadedAt:   file.UploadedAt,
		UploadedBy:   userInfo,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}
}

// getUserInfo fetches user info with username, email, and role
func getUserInfo(db *mongo.Database, userID string) model.UserInfo {
	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		return model.UserInfo{
			Username: "Unknown",
			Email:    "unknown@example.com",
			Role:     "unknown",
		}
	}

	return model.UserInfo{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}
