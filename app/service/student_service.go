package service

import (
	"context"
	"errors"

	pgModel "clean-arch-copy/app/model/postgre"
	pgRepo "clean-arch-copy/app/repository/postgre"

	"github.com/google/uuid"
)

type StudentService struct {
	repo pgRepo.StudentRepository
}

func NewStudentService(r pgRepo.StudentRepository) *StudentService {
	return &StudentService{repo: r}
}

func (s *StudentService) Create(ctx context.Context, st *pgModel.Student) error {
	if st.UserID == "" || st.StudentID == "" {
		return errors.New("missing required fields")
	}
	if st.ID == "" {
		st.ID = uuid.New().String()
	}
	return s.repo.Create(ctx, st)
}

func (s *StudentService) GetByUserID(ctx context.Context, userID string) (*pgModel.Student, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *StudentService) ListByAdvisor(ctx context.Context, advisorID string) ([]*pgModel.Student, error) {
	return s.repo.ListByAdvisor(ctx, advisorID)
}
