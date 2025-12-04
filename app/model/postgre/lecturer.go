package postgres

import "time"

type Lecturer struct {
	ID         string    `db:"id" json:"id"`                   // uuid
	UserID     string    `db:"user_id" json:"user_id"`         // FK -> users.id
	LecturerID string    `db:"lecturer_id" json:"lecturer_id"` // kode dosen
	Department string    `db:"department" json:"department"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
