package postgre

import (
	"context"
	"database/sql"
	"time"

	pgmodel "clean-arch-copy/app/model/postgre"
)

// AchievementRefRepository handles achievement_references table.
type AchievementRefRepository interface {
	Create(ctx context.Context, ref *pgmodel.AchievementReference) error
	UpdateStatus(ctx context.Context, id string, status string, verifierID *string) error
	GetByID(ctx context.Context, id string) (*pgmodel.AchievementReference, error)
	ListByStudent(ctx context.Context, studentID string) ([]*pgmodel.AchievementReference, error)
	// Extra helper: UpdateRejectionNote if needed
	UpdateRejectionNote(ctx context.Context, id string, note string) error
}

// Implementation
type achievementRefRepository struct {
	db *sql.DB
}

func NewAchievementRefRepository(db *sql.DB) AchievementRefRepository {
	return &achievementRefRepository{db: db}
}

func (r *achievementRefRepository) Create(ctx context.Context, ref *pgmodel.AchievementReference) error {
	now := time.Now()
	ref.CreatedAt = now
	ref.UpdatedAt = now
	q := `INSERT INTO achievement_references
	      (id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.ExecContext(ctx, q,
		ref.ID, ref.StudentID, ref.MongoAchievementID, ref.Status,
		ref.SubmittedAt, ref.VerifiedAt, ref.VerifiedBy, ref.RejectionNote,
		ref.CreatedAt, ref.UpdatedAt,
	)
	return err
}

func (r *achievementRefRepository) UpdateStatus(ctx context.Context, id string, status string, verifierID *string) error {
	now := time.Now()
	if verifierID != nil {
		q := `UPDATE achievement_references SET status=$1, verified_by=$2, verified_at=$3, updated_at=$4 WHERE id=$5`
		_, err := r.db.ExecContext(ctx, q, status, *verifierID, now, now, id)
		return err
	}
	q := `UPDATE achievement_references SET status=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.ExecContext(ctx, q, status, now, id)
	return err
}

func (r *achievementRefRepository) GetByID(ctx context.Context, id string) (*pgmodel.AchievementReference, error) {
	var out pgmodel.AchievementReference
	q := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	      FROM achievement_references WHERE id=$1`
	row := r.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(&out.ID, &out.StudentID, &out.MongoAchievementID, &out.Status,
		&out.SubmittedAt, &out.VerifiedAt, &out.VerifiedBy, &out.RejectionNote, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *achievementRefRepository) ListByStudent(ctx context.Context, studentID string) ([]*pgmodel.AchievementReference, error) {
	q := `SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	      FROM achievement_references WHERE student_id=$1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*pgmodel.AchievementReference
	for rows.Next() {
		var item pgmodel.AchievementReference
		if err := rows.Scan(&item.ID, &item.StudentID, &item.MongoAchievementID, &item.Status,
			&item.SubmittedAt, &item.VerifiedAt, &item.VerifiedBy, &item.RejectionNote, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &item)
	}
	return out, nil
}

func (r *achievementRefRepository) UpdateRejectionNote(ctx context.Context, id string, note string) error {
	now := time.Now()
	q := `UPDATE achievement_references SET rejection_note=$1, status='rejected', updated_at=$2 WHERE id=$3`
	_, err := r.db.ExecContext(ctx, q, note, now, id)
	return err
}
