package service

import (
	"context"
	"errors"
	"time"

	mongoModel "clean-arch-copy/app/model/mongo"
	pgModel "clean-arch-copy/app/model/postgre"
	mongoRepo "clean-arch-copy/app/repository/mongo"
	pgRepo "clean-arch-copy/app/repository/postgre"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AchievementService orchestrates Mongo + Postgres for achievements, and writes activity logs.
type AchievementService struct {
	achievementMongo mongoRepo.AchievementRepository
	achievementRefPG pgRepo.AchievementRefRepository
	studentRepo      pgRepo.StudentRepository
	userRepo         pgRepo.UserRepository
	activityRepo     pgRepo.ActivityLogRepository
}

// NewAchievementService creates an instance of AchievementService.
// NOTE: activityRepo can be nil if you don't want logging (but recommended to provide).
func NewAchievementService(
	achievementMongo mongoRepo.AchievementRepository,
	achievementRefPG pgRepo.AchievementRefRepository,
	studentRepo pgRepo.StudentRepository,
	userRepo pgRepo.UserRepository,
	activityRepo pgRepo.ActivityLogRepository,
) *AchievementService {
	return &AchievementService{
		achievementMongo: achievementMongo,
		achievementRefPG: achievementRefPG,
		studentRepo:      studentRepo,
		userRepo:         userRepo,
		activityRepo:     activityRepo,
	}
}

// helper: create activity log best-effort
func (s *AchievementService) writeActivityLog(ctx context.Context, logEntry *pgModel.ActivityLog) {
	if s.activityRepo == nil || logEntry == nil {
		return
	}
	// do not propagate error to caller; best-effort
	_ = s.activityRepo.Create(ctx, logEntry)
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
		// cleanup mongo doc best-effort
		_ = s.achievementMongo.SoftDelete(ctx, oid)
		return nil, err
	}

	// 4. write activity log (created)
	logEntry := &pgModel.ActivityLog{
		ID:         uuid.New().String(),
		EntityType: "achievement_reference",
		EntityID:   ref.ID,
		EventType:  "created",
		ActorID:    &userID,
		Previous:   nil,
		Current: map[string]interface{}{
			"status":               ref.Status,
			"mongo_achievement_id": ref.MongoAchievementID,
		},
		CreatedAt: time.Now(),
	}
	s.writeActivityLog(ctx, logEntry)

	return ref, nil
}

// Submit transitions draft -> submitted
func (s *AchievementService) Submit(ctx context.Context, refID string, userID string) error {
	// validate student
	student, err := s.studentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// get reference
	ref, err := s.achievementRefPG.GetByID(ctx, refID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement reference not found")
	}
	if ref.StudentID != student.ID {
		return errors.New("not owner")
	}
	if ref.Status != "draft" {
		return errors.New("invalid status transition: only draft can be submitted")
	}

	// update status
	now := time.Now()
	ref.SubmittedAt = &now
	if err := s.achievementRefPG.UpdateStatus(ctx, refID, "submitted", nil); err != nil {
		return err
	}

	// activity log
	logEntry := &pgModel.ActivityLog{
		ID:         uuid.New().String(),
		EntityType: "achievement_reference",
		EntityID:   ref.ID,
		EventType:  "status_changed",
		ActorID:    &userID,
		Previous:   map[string]interface{}{"status": "draft"},
		Current:    map[string]interface{}{"status": "submitted", "submitted_at": now},
		CreatedAt:  time.Now(),
	}
	s.writeActivityLog(ctx, logEntry)
	return nil
}

// Verify transitions submitted -> verified
func (s *AchievementService) Verify(ctx context.Context, refID string, verifierUserID string) error {
	// verifier existence check
	verifier, err := s.userRepo.GetByID(ctx, verifierUserID)
	if err != nil {
		return err
	}
	if verifier == nil {
		return errors.New("verifier user not found")
	}

	// get reference
	ref, err := s.achievementRefPG.GetByID(ctx, refID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement reference not found")
	}
	if ref.Status != "submitted" {
		return errors.New("only submitted achievements can be verified")
	}

	// update status in db (use UpdateStatus which sets verified_by & verified_at when provided)
	if err := s.achievementRefPG.UpdateStatus(ctx, refID, "verified", &verifierUserID); err != nil {
		return err
	}

	// activity log
	now := time.Now()
	logEntry := &pgModel.ActivityLog{
		ID:         uuid.New().String(),
		EntityType: "achievement_reference",
		EntityID:   ref.ID,
		EventType:  "status_changed",
		ActorID:    &verifierUserID,
		ActorRole:  nil, // optional: you can fetch role name if needed
		Previous:   map[string]interface{}{"status": "submitted"},
		Current:    map[string]interface{}{"status": "verified", "verified_at": now, "verified_by": verifierUserID},
		CreatedAt:  time.Now(),
	}
	s.writeActivityLog(ctx, logEntry)
	return nil
}

// Reject sets status to rejected and saves rejection note
func (s *AchievementService) Reject(ctx context.Context, refID string, verifierUserID string, note string) error {
	// verifier existence check
	verifier, err := s.userRepo.GetByID(ctx, verifierUserID)
	if err != nil {
		return err
	}
	if verifier == nil {
		return errors.New("verifier user not found")
	}

	// get reference
	ref, err := s.achievementRefPG.GetByID(ctx, refID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement reference not found")
	}
	if ref.Status != "submitted" {
		return errors.New("only submitted achievements can be rejected")
	}

	// update rejection note and status
	if err := s.achievementRefPG.UpdateRejectionNote(ctx, refID, note); err != nil {
		return err
	}

	// activity log
	now := time.Now()
	logEntry := &pgModel.ActivityLog{
		ID:         uuid.New().String(),
		EntityType: "achievement_reference",
		EntityID:   ref.ID,
		EventType:  "status_changed",
		ActorID:    &verifierUserID,
		Previous:   map[string]interface{}{"status": "submitted"},
		Current:    map[string]interface{}{"status": "rejected", "rejection_note": note, "rejected_at": now},
		CreatedAt:  time.Now(),
	}
	s.writeActivityLog(ctx, logEntry)
	return nil
}

// DeleteDraft: soft delete in Mongo + update reference in Postgres to 'deleted' (only owner, only draft)
func (s *AchievementService) DeleteDraft(ctx context.Context, refID string, userID string) error {
	// validate student
	student, err := s.studentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// get reference
	ref, err := s.achievementRefPG.GetByID(ctx, refID)
	if err != nil {
		return err
	}
	if ref == nil {
		return errors.New("achievement reference not found")
	}
	if ref.StudentID != student.ID {
		return errors.New("not owner")
	}
	if ref.Status != "draft" {
		return errors.New("only draft achievements can be deleted")
	}

	// soft delete mongo doc
	oid, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	if err != nil {
		// still update postgres to deleted for safety
		_ = s.achievementRefPG.UpdateStatus(ctx, refID, "deleted", nil)
		return errors.New("invalid mongo object id")
	}
	if err := s.achievementMongo.SoftDelete(ctx, oid); err != nil {
		// continue to update postgres anyway (best-effort)
	}

	// update postgres reference status to deleted
	if err := s.achievementRefPG.UpdateStatus(ctx, refID, "deleted", nil); err != nil {
		return err
	}

	// activity log
	now := time.Now()
	logEntry := &pgModel.ActivityLog{
		ID:         uuid.New().String(),
		EntityType: "achievement_reference",
		EntityID:   ref.ID,
		EventType:  "deleted",
		ActorID:    &userID,
		Previous:   map[string]interface{}{"status": "draft"},
		Current:    map[string]interface{}{"status": "deleted", "deleted_at": now},
		CreatedAt:  time.Now(),
	}
	s.writeActivityLog(ctx, logEntry)
	return nil
}

// GetDetail returns both Mongo document and Postgres reference
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
