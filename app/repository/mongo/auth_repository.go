package repository

import (
	"clean-arch/app/model/mongo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByUsernameOrEmail(db *mongo.Database, identifier string) (*model.User, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection("users")

	var result bson.M
	err := collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": identifier},
			{"email": identifier},
		},
	}).Decode(&result)

	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		ID:       result["_id"].(primitive.ObjectID),
		Username: result["username"].(string),
		Email:    result["email"].(string),
		Role:     result["role"].(string),
	}

	passwordHash, _ := result["password_hash"].(string)

	return user, passwordHash, nil
}

func GetAlumniByNIM(db *mongo.Database, nim string) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	var alumni model.Alumni
	err := collection.FindOne(ctx, bson.M{"nim": nim}).Decode(&alumni)
	if err != nil {
		return nil, err
	}

	return &alumni, nil
}

func GetAlumniWithJobs(db *mongo.Database, alumniID string) (*model.AlumniWithJobs, error) {
	alumni, err := GetAlumniByID(db, alumniID)
	if err != nil {
		return nil, err
	}

	jobs, err := GetPekerjaanByAlumniID(db, alumniID)
	if err != nil {
		return nil, err
	}

	return &model.AlumniWithJobs{
		Alumni:        *alumni,
		PekerjaanList: jobs,
	}, nil
}

func CreateAlumniWithAuth(db *mongo.Database, req model.CreateAlumniRequest, hashedPassword string) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)
	now := time.Now()

	role := "user"
	if req.UserID != nil {
		role = "admin"
	}

	alumni := model.Alumni{
		ID:         primitive.NewObjectID(),
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       role,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result, err := collection.InsertOne(ctx, alumni)
	if err != nil {
		return nil, err
	}

	alumni.ID = result.InsertedID.(primitive.ObjectID)
	return &alumni, nil
}

func GetUserByID(db *mongo.Database, userID string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
