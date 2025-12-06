package postgre

import (
	"context"
	"database/sql"
	"time"

	pgmodel "clean-arch-copy/app/model/postgre"
)

// LecturerRepository manages lecturers table.
type LecturerRepository interface {
	Create(ctx context.Context, l *pgmodel.Lecturer) error
	GetByID(ctx context.Context, id string) (*pgmodel.Lecturer, error)
	GetByUserID(ctx context.Context, userID string) (*pgmodel.Lecturer, error)
	ListAll(ctx context.Context) ([]*pgmodel.Lecturer, error)
}

// Implementation
type lecturerRepository struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepository{db: db}
}

func (r *lecturerRepository) Create(ctx context.Context, l *pgmodel.Lecturer) error {
	l.CreatedAt = time.Now()
	q := `INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.ExecContext(ctx, q, l.ID, l.UserID, l.LecturerID, l.Department, l.CreatedAt)
	return err
}

func (r *lecturerRepository) GetByID(ctx context.Context, id string) (*pgmodel.Lecturer, error) {
	var out pgmodel.Lecturer
	q := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE id=$1`
	row := r.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(&out.ID, &out.UserID, &out.LecturerID, &out.Department, &out.CreatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *lecturerRepository) GetByUserID(ctx context.Context, userID string) (*pgmodel.Lecturer, error) {
	var out pgmodel.Lecturer
	q := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE user_id=$1`
	row := r.db.QueryRowContext(ctx, q, userID)
	if err := row.Scan(&out.ID, &out.UserID, &out.LecturerID, &out.Department, &out.CreatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *lecturerRepository) ListAll(ctx context.Context) ([]*pgmodel.Lecturer, error) {
	q := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers ORDER BY lecturer_id`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*pgmodel.Lecturer
	for rows.Next() {
		var l pgmodel.Lecturer
		if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &l)
	}
	return out, nil
}
