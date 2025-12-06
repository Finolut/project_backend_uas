package mongo

import (
	"context"
	"errors"
	"time"

	mongomodel "clean-arch-copy/app/model/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AchievementRepository defines the operations for achievements collection.
type AchievementRepository interface {
	Create(ctx context.Context, a *mongomodel.Achievement) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*mongomodel.Achievement, error)
	Update(ctx context.Context, a *mongomodel.Achievement) error
	SoftDelete(ctx context.Context, id primitive.ObjectID) error

	// Optional helpers
	ListByStudent(ctx context.Context, studentID string, limit, offset int64) ([]*mongomodel.Achievement, error)
}

// --------------------------
// Implementation
// --------------------------
type achievementRepo struct {
	col *driver.Collection
}

// NewAchievementRepository creates repository and ensures indexes
func NewAchievementRepository(db *driver.Database, collectionName string) AchievementRepository {
	col := db.Collection(collectionName)
	r := &achievementRepo{col: col}

	// ensure indexes (best-effort)
	_ = r.ensureIndexes(context.Background())
	return r
}

func (r *achievementRepo) ensureIndexes(ctx context.Context) error {
	// index on studentId and createdAt
	indexes := []driver.IndexModel{
		{
			Keys:    bson.D{{Key: "studentId", Value: 1}},
			Options: options.Index().SetBackground(true),
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetBackground(true),
		},
		{
			Keys:    bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetBackground(true),
		},
	}
	_, err := r.col.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create inserts a new achievement document and returns its ObjectID
func (r *achievementRepo) Create(ctx context.Context, a *mongomodel.Achievement) (primitive.ObjectID, error) {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	res, err := r.col.InsertOne(ctx, a)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("inserted id is not an ObjectID")
	}
	return oid, nil
}

// GetByID fetches an achievement by ObjectID
func (r *achievementRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*mongomodel.Achievement, error) {
	var out mongomodel.Achievement
	filter := bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}
	err := r.col.FindOne(ctx, filter).Decode(&out)
	if err != nil {
		if err == driver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

// Update replaces the whole document (caller must ensure ID present)
func (r *achievementRepo) Update(ctx context.Context, a *mongomodel.Achievement) error {
	if a.ID.IsZero() {
		return errors.New("missing achievement ID")
	}
	a.UpdatedAt = time.Now()
	_, err := r.col.ReplaceOne(ctx, bson.M{"_id": a.ID}, a)
	return err
}

// SoftDelete sets deletedAt timestamp instead of physically deleting
func (r *achievementRepo) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"deletedAt": now,
			"updatedAt": now,
		},
	}
	res, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return driver.ErrNoDocuments
	}
	return nil
}

// ListByStudent returns achievements for a given student with pagination
func (r *achievementRepo) ListByStudent(ctx context.Context, studentID string, limit, offset int64) ([]*mongomodel.Achievement, error) {
	if limit <= 0 {
		limit = 20
	}
	filter := bson.M{
		"studentId": studentID,
		"deletedAt": bson.M{"$exists": false},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(limit).
		SetSkip(offset)

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []*mongomodel.Achievement
	for cur.Next(ctx) {
		var a mongomodel.Achievement
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		out = append(out, &a)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
