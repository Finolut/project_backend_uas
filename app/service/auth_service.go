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

func (s *AuthService) Refresh(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user_id is required")
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(s.jwtSecret))
	return ss, err
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}
	// TODO: implement token blacklist or revocation mechanism
	// For now, this is a placeholder for future implementation
	return nil
}

func (s *AuthService) VerifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}
