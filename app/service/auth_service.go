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

// Definisikan interface untuk Token Repository di sini (atau import dari domain layer)
// Implementasinya nanti bisa menggunakan Redis (disarankan) atau Database SQL
type TokenRepository interface {
	AddToBlacklist(ctx context.Context, token string, expiresAt time.Time) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type AuthService struct {
	userRepo  pgRepo.UserRepository
	tokenRepo TokenRepository // Tambahkan dependency ini
	jwtSecret string
}

// Update constructor untuk menerima tokenRepo
// Note: Anda perlu mengupdate wiring di service_factory.go juga nantinya
func NewAuthService(userRepo pgRepo.UserRepository, tokenRepo TokenRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: secret,
	}
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

// Logout memasukkan token ke dalam blacklist hingga masa berlakunya habis
// NOTE: Saya mengubah parameter userID menjadi tokenString karena untuk blacklist kita butuh tokennya
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	if tokenString == "" {
		return errors.New("token is required")
	}

	// 1. Parse token tanpa verifikasi signature (hanya butuh klaim 'exp')
	// Kita bisa gunakan ParseUnverified karena tujuan kita hanya membaca expiration time
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// 2. Ambil waktu kadaluarsa (exp)
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("token does not have expiration time")
	}
	
	expiresAt := time.Unix(int64(expFloat), 0)
	
	// Jika token sudah expired, tidak perlu di-blacklist
	if time.Now().After(expiresAt) {
		return nil 
	}

	// 3. Simpan token ke blacklist repository
	// Token akan disimpan di DB/Redis sampai waktu 'expiresAt' tercapai
	return s.tokenRepo.AddToBlacklist(ctx, tokenString, expiresAt)
}

func (s *AuthService) VerifyToken(tokenString string) (map[string]interface{}, error) {
	// 1. Cek apakah token ada di blacklist
	if s.tokenRepo != nil {
		isBlacklisted, err := s.tokenRepo.IsBlacklisted(context.Background(), tokenString)
		if err != nil {
			// Fail safe: jika DB error, anggap invalid atau log error
			return nil, err 
		}
		if isBlacklisted {
			return nil, errors.New("token has been invalidated (logged out)")
		}
	}

	// 2. Standard JWT verification
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