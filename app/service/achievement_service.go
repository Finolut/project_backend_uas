package service

import (
	"context"
	"errors"
	"time"

	mongoModel "app/model/mongo"
	pgModel "app/model/postgre"
	mongoRepo "app/repository/mongo"
	pgRepo "app/repository/postgre"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService struct {
	achievementMongo mongoRepo.AchievementRepository
	achievementRefPG pgRepo.AchievementRefRepository
	studentRepo      pgRepo.StudentRepository
	userRepo         pgRepo.UserRepository
}

func NewAchievementService(
	achievementMongo mongoRepo.AchievementRepository,
	achievementRefPG pgRepo.AchievementRefRepository,
	studentRepo pgRepo.StudentRepository,
	userRepo pgRepo.UserRepository,
) *AchievementService {
	return &AchievementService{
		achievementMongo: achievementMongo,
		achievementRefPG: achievementRefPG,
		studentRepo:      studentRepo,
		userRepo:         userRepo,
	}
}

// CreateDraft saves achievement doc to Mongo and creates a reference row in Postgres (status=draft)
func (s *AchievementService) CreateDraft(ctx context.Context, userID string, doc *mongoModel.Achievement) (*pgModel.AchievementReference, error) {
	// 1. validate student
	student, err := s.studentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("student profile not found")
	}

	// 2. save to mongo
	doc.StudentID = student.ID
	oid, err := s.achievementMongo.Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	// 3. create reference in postgres
	ref := &pgModel.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          student.ID,
		MongoAchievementID: oid.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := s.achievementRefPG.Create(ctx, ref); err != nil {
		// try to cleanup mongo doc (best-effort)
		_ = s.achievementMongo.SoftDelete(ctx, oid)
		return nil, err
	}
	return ref, nil
}

func (s *AchievementService) Submit(ctx context.Context, refID string, userID string) error {
	student, err := s.studentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}
	ref, err := s.achievementRefPG.GetByID(ctx, refID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("reference not found")
	}
	if ref.StudentID != student.ID {
		return errors.New("not owner")
	}
	if ref.Status != "draft" {
		return errors.New("invalid status transition")
	}

	now := time.Now()
	ref.SubmittedAt = &now
	return s.achievementRefPG.UpdateStatus(ctx, refID, "submitted", nil)
}

func (s *AchievementService) Verify(ctx context.Context, refID string, verifierUserID string) error {
	// verifier existence check
	verifier, err := s.userRepo.GetByID(ctx, verifierUserID)
	if err != nil {
		return err
	}
	if verifier == nil {
		return errors.New("verifier not found")
	}

	ref, err := s.achievementRefPG.GetByID(ctx, refID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("reference not found")
	}
	if ref.Status != "submitted" {
		return errors.New("only submitted achievements can be verified")
	}

	return s.achievementRefPG.UpdateStatus(ctx, refID, "verified", &verifierUserID)
}

func (s *AchievementService) Reject(ctx context.Context, refID string, verifierUserID string, note string) error {
	verifier, err := s.userRepo.GetByID(ctx, verifierUserID)
	if err != nil {
		return err
	}
	if verifier == nil {
		return errors.New("verifier not found")
	}
	// mark rejected with note
	return s.achievementRefPG.UpdateRejectionNote(ctx, refID, note)
}

func (s *AchievementService) GetDetail(ctx context.Context, refID string) (*mongoModel.Achievement, *pgModel.AchievementReference, error) {
	ref, err := s.achievementRefPG.GetByID(ctx, refID)
	if err != nil {
		return nil, nil, err
	}
	if ref == nil {
		return nil, nil, errors.New("reference not found")
	}
	oid, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		return nil, nil, errors.New("invalid mongo id stored in reference")
	}
	ach, err := s.achievementMongo.GetByID(ctx, oid)
	if err != nil {
		return nil, nil, err
	}
	return ach, ref, nil
}
