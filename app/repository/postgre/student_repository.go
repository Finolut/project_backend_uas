package postgre

import (
	"context"
	"database/sql"
	"time"

	pgmodel "app/model/postgre"
)

// StudentRepository provides CRUD for students.
type StudentRepository interface {
	Create(ctx context.Context, s *pgmodel.Student) error
	GetByID(ctx context.Context, id string) (*pgmodel.Student, error)
	GetByUserID(ctx context.Context, userID string) (*pgmodel.Student, error)
	ListByAdvisor(ctx context.Context, advisorID string) ([]*pgmodel.Student, error)
}

// Implementation
type studentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, s *pgmodel.Student) error {
	s.CreatedAt = time.Now()
	q := `INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
	      VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.ExecContext(ctx, q, s.ID, s.UserID, s.StudentID, s.Program, s.AcademicYear, s.AdvisorID, s.CreatedAt)
	return err
}

func (r *studentRepository) GetByID(ctx context.Context, id string) (*pgmodel.Student, error) {
	var out pgmodel.Student
	q := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE id=$1`
	row := r.db.QueryRowContext(ctx, q, id)
	if err := row.Scan(&out.ID, &out.UserID, &out.StudentID, &out.Program, &out.AcademicYear, &out.AdvisorID, &out.CreatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *studentRepository) GetByUserID(ctx context.Context, userID string) (*pgmodel.Student, error) {
	var out pgmodel.Student
	q := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE user_id=$1`
	row := r.db.QueryRowContext(ctx, q, userID)
	if err := row.Scan(&out.ID, &out.UserID, &out.StudentID, &out.Program, &out.AcademicYear, &out.AdvisorID, &out.CreatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *studentRepository) ListByAdvisor(ctx context.Context, advisorID string) ([]*pgmodel.Student, error) {
	q := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE advisor_id=$1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []*pgmodel.Student{}
	for rows.Next() {
		var s pgmodel.Student
		if err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.Program, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}
	return out, nil
}
