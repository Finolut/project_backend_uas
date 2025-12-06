package service

import (
	"context"
	"errors"
	"os"
	"time"

	pgModel "clean-arch-copy/app/model/postgre"
	pgRepo "clean-arch-copy/app/repository/postgre"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo pgRepo.UserRepository
	// optionally: token store, refresh token repo
	jwtSecret string
}

func NewAuthService(userRepo pgRepo.UserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret" // TODO: fail fast in production
	}
	return &AuthService{userRepo: userRepo, jwtSecret: secret}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func (s *AuthService) ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Login authenticates and returns JWT token
func (s *AuthService) Login(ctx context.Context, username, password string) (string, *pgModel.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := s.ComparePassword(user.PasswordHash, password); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	// create token
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.RoleID,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return ss, user, nil
}
