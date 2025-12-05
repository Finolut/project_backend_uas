package service

import (
	"context"
	"errors"

	pgModel "app/model/postgre"
	pgRepo "app/repository/postgre"

	"github.com/google/uuid"
)

type LecturerService struct {
	repo pgRepo.LecturerRepository
}

func NewLecturerService(r pgRepo.LecturerRepository) *LecturerService {
	return &LecturerService{repo: r}
}

func (s *LecturerService) Create(ctx context.Context, l *pgModel.Lecturer) error {
	if l.UserID == "" || l.LecturerID == "" {
		return errors.New("missing required fields")
	}
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return s.repo.Create(ctx, l)
}

func (s *LecturerService) GetByUserID(ctx context.Context, userID string) (*pgModel.Lecturer, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *LecturerService) ListAll(ctx context.Context) ([]*pgModel.Lecturer, error) {
	return s.repo.ListAll(ctx)
}
