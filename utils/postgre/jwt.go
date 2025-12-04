package utils

import (
	model "clean-arch/app/model/postgre"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-min-32-characters-long")
var refreshTokenSecret = []byte("your-refresh-secret-key-min-32-characters-long")

// GenerateToken is updated to use new User model and include RoleName
func GenerateToken(user model.User, roleName string) (string, error) {
	claims := model.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleName: roleName,
		RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(user model.User, roleName string) (string, error) {
	claims := model.RefreshJWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleName: roleName,
		RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenSecret)
}

func GenerateStudentToken(user model.User, student model.Student, roleName string) (string, error) {
	claims := model.StudentJWTClaims{
		UserID:   user.ID,
		NIM:      student.StudentID,
		Nama:     user.FullName,
		Email:    user.Email,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

func ValidateRefreshToken(tokenString string) (*model.RefreshJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.RefreshJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return refreshTokenSecret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.RefreshJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

func ValidateStudentToken(tokenString string) (*model.StudentJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.StudentJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.StudentJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
