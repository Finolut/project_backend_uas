package postgres

import "time"

type AchievementReference struct {
	ID                 string     `db:"id" json:"id"`                                     // uuid
	StudentID          string     `db:"student_id" json:"student_id"`                     // FK -> students.id
	MongoAchievementID string     `db:"mongo_achievement_id" json:"mongo_achievement_id"` // ObjectId.Hex()
	Status             string     `db:"status" json:"status"`                             // draft, submitted, verified, rejected
	SubmittedAt        *time.Time `db:"submitted_at" json:"submitted_at"`
	VerifiedAt         *time.Time `db:"verified_at" json:"verified_at"`
	VerifiedBy         *string    `db:"verified_by" json:"verified_by"` // FK -> users.id (verifier)
	RejectionNote      *string    `db:"rejection_note" json:"rejection_note"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}
