package repository

import (
	"clean-arch/app/model/mongo"
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const alumniCollection = "alumni"

func GetAllAlumniWithPagination(db *mongo.Database, params model.PaginationParams) ([]model.Alumni, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	// Build filter for search
	filter := bson.M{"deleted_at": nil}
	if params.Search != "" {
		filter = bson.M{
			"$and": []bson.M{
				{"deleted_at": nil},
				{
					"$or": []bson.M{
						{"nama": bson.M{"$regex": params.Search, "$options": "i"}},
						{"nim": bson.M{"$regex": params.Search, "$options": "i"}},
						{"jurusan": bson.M{"$regex": params.Search, "$options": "i"}},
						{"email": bson.M{"$regex": params.Search, "$options": "i"}},
					},
				},
			},
		}
	}

	// Get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Validate sort column
	validSortColumns := map[string]bool{
		"nim": true, "nama": true, "jurusan": true,
		"angkatan": true, "tahun_lulus": true, "email": true, "created_at": true,
	}
	if !validSortColumns[params.SortBy] {
		params.SortBy = "created_at"
	}

	// Validate sort order
	sortOrder := 1
	if strings.ToLower(params.Order) == "desc" {
		sortOrder = -1
	}

	// Query with pagination
	offset := int64((params.Page - 1) * params.Limit)
	opts := options.Find().
		SetSort(bson.M{params.SortBy: sortOrder}).
		SetSkip(offset).
		SetLimit(int64(params.Limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var alumniList []model.Alumni
	if err = cursor.All(ctx, &alumniList); err != nil {
		return nil, 0, err
	}

	return alumniList, int(total), nil
}

func GetAllAlumni(db *mongo.Database) ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)
	filter := bson.M{"deleted_at": nil}
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumniList []model.Alumni
	if err = cursor.All(ctx, &alumniList); err != nil {
		return nil, err
	}

	return alumniList, nil
}

func GetAlumniByID(db *mongo.Database, id string) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var alumni model.Alumni
	err = collection.FindOne(ctx, bson.M{"_id": objID, "deleted_at": nil}).Decode(&alumni)
	if err != nil {
		return nil, err
	}

	return &alumni, nil
}

func CreateAlumni(db *mongo.Database, req model.CreateAlumniRequest) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)
	now := time.Now()

	alumni := model.Alumni{
		ID:         primitive.NewObjectID(),
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		Password:   req.Password,
		Role:       "user",
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

func UpdateAlumni(db *mongo.Database, id string, req model.UpdateAlumniRequest) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"nama":        req.Nama,
			"jurusan":     req.Jurusan,
			"angkatan":    req.Angkatan,
			"tahun_lulus": req.TahunLulus,
			"email":       req.Email,
			"no_telepon":  req.NoTelepon,
			"alamat":      req.Alamat,
			"updated_at":  time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedAlumni model.Alumni

	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": objID, "deleted_at": nil}, update, opts).Decode(&updatedAlumni)
	if err != nil {
		return nil, err
	}

	return &updatedAlumni, nil
}

func DeleteAlumni(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func CheckAlumniByNim(db *mongo.Database, nim string) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	var alumni model.Alumni
	err := collection.FindOne(ctx, bson.M{"nim": nim, "deleted_at": nil}).Decode(&alumni)
	if err != nil {
		return nil, err
	}

	return &alumni, nil
}

func GetAlumniStatistics(db *mongo.Database) (*model.AlumniStatistics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)
	filter := bson.M{"deleted_at": nil}

	stats := &model.AlumniStatistics{
		AlumniByJurusan:    make(map[string]int),
		AlumniByAngkatan:   make(map[string]int),
		AlumniByTahunLulus: make(map[string]int),
	}

	// Get total alumni count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	stats.TotalAlumni = int(total)

	// Get alumni count by jurusan
	cursor, err := collection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "$jurusan", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jurusanResults []bson.M
	if err = cursor.All(ctx, &jurusanResults); err != nil {
		return nil, err
	}

	for _, result := range jurusanResults {
		jurusan := result["_id"].(string)
		count := int(result["count"].(int32))
		stats.AlumniByJurusan[jurusan] = count
	}

	// Similar aggregations for angkatan and tahun_lulus
	angkatanCursor, err := collection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "$angkatan", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	defer angkatanCursor.Close(ctx)

	var angkatanResults []bson.M
	if err = angkatanCursor.All(ctx, &angkatanResults); err != nil {
		return nil, err
	}

	for _, result := range angkatanResults {
		angkatan := fmt.Sprintf("%v", result["_id"])
		count := int(result["count"].(int32))
		stats.AlumniByAngkatan[angkatan] = count
	}

	tahunCursor, err := collection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "$tahun_lulus", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	defer tahunCursor.Close(ctx)

	var tahunResults []bson.M
	if err = tahunCursor.All(ctx, &tahunResults); err != nil {
		return nil, err
	}

	for _, result := range tahunResults {
		tahun := fmt.Sprintf("%v", result["_id"])
		count := int(result["count"].(int32))
		stats.AlumniByTahunLulus[tahun] = count
	}

	return stats, nil
}

func GetTrashedAlumni(db *mongo.Database) ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)
	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	opts := options.Find().SetSort(bson.M{"deleted_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Alumni
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func SoftDeleteAlumni(db *mongo.Database, id string, userID *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"deleted_by": userID,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID, "deleted_at": nil}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func RestoreAlumni(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": nil,
			"deleted_by": nil,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID, "deleted_at": bson.M{"$ne": nil}}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func HardDeleteAlumni(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(alumniCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
