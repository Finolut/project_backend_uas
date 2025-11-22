package repository

import (
	"clean-arch/app/model/mongo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const fileCollection = "files"

// CreateFile saves file metadata to database
func CreateFile(db *mongo.Database, file *model.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(fileCollection)
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}

	file.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetFileByUserID retrieves files for a specific user
func GetFileByUserID(db *mongo.Database, userID string, category string) ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(fileCollection)

	filter := bson.M{
		"user_id":    userID,
		"category":   category,
		"deleted_at": nil,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []model.File
	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	return files, nil
}

// GetFileByID retrieves a specific file by ID
func GetFileByID(db *mongo.Database, id string) (*model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(fileCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var file model.File
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "deleted_at": nil}).Decode(&file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

// DeleteFile performs soft delete on file
func DeleteFile(db *mongo.Database, id string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(fileCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$set": bson.M{
			"deleted_at": now,
			"updated_at": now,
		},
	})
	return err
}

// GetAllFilesByCategory retrieves all files of a specific category (admin only)
func GetAllFilesByCategory(db *mongo.Database, category string) ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(fileCollection)

	filter := bson.M{
		"category":   category,
		"deleted_at": nil,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []model.File
	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	return files, nil
}
