package repository

import (
	"clean-arch/app/model"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const pekerjaanCollection = "pekerjaan_alumni"

func GetAllPekerjaanWithPagination(db *mongo.Database, params model.PaginationParams) ([]model.PekerjaanAlumni, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	// Build filter
	filter := bson.M{"deleted_at": nil}
	if params.Search != "" {
		filter = bson.M{
			"$and": []bson.M{
				{"deleted_at": nil},
				{
					"$or": []bson.M{
						{"nama_perusahaan": bson.M{"$regex": params.Search, "$options": "i"}},
						{"posisi_jabatan": bson.M{"$regex": params.Search, "$options": "i"}},
						{"bidang_industri": bson.M{"$regex": params.Search, "$options": "i"}},
						{"lokasi_kerja": bson.M{"$regex": params.Search, "$options": "i"}},
						{"status_pekerjaan": bson.M{"$regex": params.Search, "$options": "i"}},
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
		"nama_perusahaan": true, "posisi_jabatan": true, "bidang_industri": true,
		"lokasi_kerja": true, "status_pekerjaan": true, "tanggal_mulai_kerja": true, "created_at": true,
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

	var pekerjaanList []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, 0, err
	}

	return pekerjaanList, int(total), nil
}

func GetAllPekerjaan(db *mongo.Database) ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)
	filter := bson.M{"deleted_at": nil}
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaanList []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, err
	}

	return pekerjaanList, nil
}

func GetPekerjaanByID(db *mongo.Database, id string) (*model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pekerjaan model.PekerjaanAlumni
	err = collection.FindOne(ctx, bson.M{"_id": objID, "deleted_at": nil}).Decode(&pekerjaan)
	if err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}

func GetPekerjaanByAlumniID(db *mongo.Database, alumniID string) ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"alumni_id": objID, "deleted_at": nil}
	opts := options.Find().SetSort(bson.M{"tanggal_mulai_kerja": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaanList []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaanList); err != nil {
		return nil, err
	}

	return pekerjaanList, nil
}

func CreatePekerjaan(db *mongo.Database, req model.CreatePekerjaanRequest, alumniID string) (*model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objAlumniID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	pekerjaan := model.PekerjaanAlumni{
		ID:                  primitive.NewObjectID(),
		AlumniID:            objAlumniID,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   req.TanggalMulaiKerja,
		TanggalSelesaiKerja: req.TanggalSelesaiKerja,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	result, err := collection.InsertOne(ctx, pekerjaan)
	if err != nil {
		return nil, err
	}

	pekerjaan.ID = result.InsertedID.(primitive.ObjectID)
	return &pekerjaan, nil
}

func UpdatePekerjaan(db *mongo.Database, id string, req model.UpdatePekerjaanRequest) (*model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan":       req.NamaPerusahaan,
			"posisi_jabatan":        req.PosisiJabatan,
			"bidang_industri":       req.BidangIndustri,
			"lokasi_kerja":          req.LokasiKerja,
			"gaji_range":            req.GajiRange,
			"tanggal_mulai_kerja":   req.TanggalMulaiKerja,
			"tanggal_selesai_kerja": req.TanggalSelesaiKerja,
			"status_pekerjaan":      req.StatusPekerjaan,
			"deskripsi_pekerjaan":   req.DeskripsiPekerjaan,
			"updated_at":            time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedPekerjaan model.PekerjaanAlumni

	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": objID, "deleted_at": nil}, update, opts).Decode(&updatedPekerjaan)
	if err != nil {
		return nil, err
	}

	return &updatedPekerjaan, nil
}

func DeletePekerjaan(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

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

func SoftDeletePekerjaan(db *mongo.Database, id string, deletedBy string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
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

func GetAlumniIDByPekerjaanID(db *mongo.Database, pekerjaanID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objID, err := primitive.ObjectIDFromHex(pekerjaanID)
	if err != nil {
		return "", err
	}

	var pekerjaan model.PekerjaanAlumni
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&pekerjaan)
	if err != nil {
		return "", err
	}

	return pekerjaan.AlumniID.Hex(), nil
}

func RestorePekerjaan(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

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

func HardDeletePekerjaan(db *mongo.Database, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

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

func SoftDeletePekerjaanByAlumniID(db *mongo.Database, alumniID string, deletedBy string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objAlumniID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateMany(ctx, bson.M{"alumni_id": objAlumniID, "deleted_at": nil}, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func HardDeletePekerjaanByAlumniID(db *mongo.Database, alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(pekerjaanCollection)

	objAlumniID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	result, err := collection.DeleteMany(ctx, bson.M{"alumni_id": objAlumniID, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
