package utils

import (
	"clean-arch/app/model/postgre"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-min-32-characters-long")

func GenerateToken(user model.User) (string, error) {
	claims := model.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateAlumniToken(alumni model.Alumni) (string, error) {
	claims := model.AlumniJWTClaims{
		AlumniID: alumni.ID,
		NIM:      alumni.NIM,
		Nama:     alumni.Nama,
		Email:    alumni.Email,
		Role:     alumni.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
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

func ValidateAlumniToken(tokenString string) (*model.AlumniJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.AlumniJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.AlumniJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
