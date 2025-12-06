package service

import (
	"context"
	"errors"
	"time"

	pgModel "clean-arch-copy/app/model/postgre"
	pgRepo "clean-arch-copy/app/repository/postgre"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo pgRepo.UserRepository
}

func NewUserService(userRepo pgRepo.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Register creates a new user (password hashing done in AuthService)
func (s *UserService) Register(ctx context.Context, u *pgModel.User) error {
	// simple validations
	if u.Username == "" || u.Email == "" || u.PasswordHash == "" {
		return errors.New("missing required fields")
	}
	u.ID = uuid.New().String()
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return s.userRepo.Create(ctx, u)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*pgModel.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*pgModel.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *UserService) Update(ctx context.Context, u *pgModel.User) error {
	u.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ListAll(ctx context.Context) ([]*pgModel.User, error) {
	return s.userRepo.ListAll(ctx)
}

func (s *UserService) UpdateRole(ctx context.Context, userID string, roleID string) error {
	if userID == "" || roleID == "" {
		return errors.New("user_id and role_id are required")
	}
	return s.userRepo.UpdateRole(ctx, userID, roleID)
}
