package service

import (
	"context"
	"errors"
	"testing"

	"clean-arch/app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -------------------- mock repository & auth helper --------------------

type AuthRepoMock interface {
	GetUserByUsernameOrEmail(ctx context.Context, username string) (*model.User, string, error)
	GetAlumniByNIM(ctx context.Context, nim string) (*model.Alumni, error)
}

type mockAuthRepo struct {
	user     *model.User
	passHash string
	alumni   *model.Alumni
	preErr   error
}

func (m *mockAuthRepo) GetUserByUsernameOrEmail(ctx context.Context, username string) (*model.User, string, error) {
	if m.preErr != nil {
		return nil, "", m.preErr
	}
	if m.user == nil {
		return nil, "", errors.New("not found")
	}
	return m.user, m.passHash, nil
}

func (m *mockAuthRepo) GetAlumniByNIM(ctx context.Context, nim string) (*model.Alumni, error) {
	if m.preErr != nil {
		return nil, m.preErr
	}
	if m.alumni == nil || m.alumni.NIM != nim {
		return nil, errors.New("not found")
	}
	return m.alumni, nil
}

// AuthHelper is a small interface to abstract password check & token generation
type AuthHelper interface {
	CheckPassword(plain, hash string) bool
	GenerateTokenForUser(u model.User) (string, error)
	GenerateTokenForAlumni(a model.Alumni) (string, error)
}

type simpleAuthHelper struct{}

func (s *simpleAuthHelper) CheckPassword(plain, hash string) bool {
	// in tests we use hash=="hashed:"+plain
	return hash == "hashed:"+plain
}

func (s *simpleAuthHelper) GenerateTokenForUser(u model.User) (string, error) {
	return "token-user-" + u.Username, nil
}

func (s *simpleAuthHelper) GenerateTokenForAlumni(a model.Alumni) (string, error) {
	return "token-alumni-" + a.NIM, nil
}

// -------------------- TestableAuthService --------------------

type TestableAuthService struct {
	repo  AuthRepoMock
	authH AuthHelper
}

func NewTestableAuthService(repo AuthRepoMock, h AuthHelper) *TestableAuthService {
	return &TestableAuthService{repo: repo, authH: h}
}

func (s *TestableAuthService) Login(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	if username == "" || password == "" {
		return nil, errors.New("username dan password harus diisi")
	}
	user, passHash, err := s.repo.GetUserByUsernameOrEmail(ctx, username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}
	if !s.authH.CheckPassword(password, passHash) {
		return nil, errors.New("username atau password salah")
	}
	tok, err := s.authH.GenerateTokenForUser(*user)
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{User: *user, Token: tok}, nil
}

func (s *TestableAuthService) AlumniLogin(ctx context.Context, nim, password string) (*model.AlumniLoginResponse, error) {
	if nim == "" || password == "" {
		return nil, errors.New("NIM dan password harus diisi")
	}
	al, err := s.repo.GetAlumniByNIM(ctx, nim)
	if err != nil {
		return nil, errors.New("NIM atau password salah")
	}
	if !s.authH.CheckPassword(password, al.Password) {
		return nil, errors.New("NIM atau password salah")
	}
	tok, err := s.authH.GenerateTokenForAlumni(*al)
	if err != nil {
		return nil, err
	}
	// remove password for response
	al.Password = ""
	return &model.AlumniLoginResponse{Alumni: *al, Token: tok}, nil
}

// -------------------- TESTS --------------------

func TestLoginService(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockAuthRepo{
		user:     &model.User{Username: "admin", Email: "admin@example.com", Role: "admin"},
		passHash: "hashed:secret123",
	}
	helper := &simpleAuthHelper{}
	svc := NewTestableAuthService(mockRepo, helper)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
		wantTok  string
	}{
		{"valid", "admin", "secret123", false, "token-user-admin"},
		{"invalid password", "admin", "wrong", true, ""},
		{"empty username", "", "x", true, ""},
		{"unknown user", "nouser", "x", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Login(ctx, tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Login error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if resp.Token != tt.wantTok {
					t.Fatalf("token = %v, want %v", resp.Token, tt.wantTok)
				}
			}
		})
	}
}

func TestAlumniLoginService(t *testing.T) {
	ctx := context.Background()

	al := &model.Alumni{
		ID:       primitive.NewObjectID(),
		NIM:      "18001",
		Nama:     "Budi",
		Email:    "budi@example.com",
		Password: "hashed:alpass",
	}
	mockRepo := &mockAuthRepo{alumni: al}
	helper := &simpleAuthHelper{}
	svc := NewTestableAuthService(mockRepo, helper)

	tests := []struct {
		name    string
		nim     string
		pass    string
		wantErr bool
		wantTok string
	}{
		{"valid", "18001", "alpass", false, "token-alumni-18001"},
		{"wrong pass", "18001", "bad", true, ""},
		{"unknown nim", "99999", "x", true, ""},
		{"empty", "", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.AlumniLogin(ctx, tt.nim, tt.pass)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AlumniLogin error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && resp.Token != tt.wantTok {
				t.Fatalf("token = %v, want %v", resp.Token, tt.wantTok)
			}
		})
	}
}
