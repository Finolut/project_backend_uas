package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"clean-arch/app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NOTE:
// File ini adalah contoh *unit tests* untuk logic service alumni dengan dependency injection.
// Karena kode asli memanggil repository package secara langsung, agar dapat dites
// Anda perlu membuat versi service yang menerima interface repository.
// Di production code Anda bisa membuat adaptor yang mengimplementasikan interface ini
// dan memanggil repository.* asli.

// -------------------- interface dan mock --------------------

type AlumniRepoMock interface {
	CheckAlumniByNim(ctx context.Context, nim string) (*model.Alumni, error)
	CreateAlumni(ctx context.Context, req model.CreateAlumniRequest) (*model.Alumni, error)
	GetAlumniByID(ctx context.Context, id string) (*model.Alumni, error)
}

type mockAlumniRepo struct {
	alumniStore map[string]*model.Alumni
	preErr      error // jika tidak nil, akan dikembalikan pada semua panggilan
}

func newMockAlumniRepo() *mockAlumniRepo {
	return &mockAlumniRepo{
		alumniStore: make(map[string]*model.Alumni),
	}
}

func (m *mockAlumniRepo) CheckAlumniByNim(ctx context.Context, nim string) (*model.Alumni, error) {
	if m.preErr != nil {
		return nil, m.preErr
	}
	for _, a := range m.alumniStore {
		if a.NIM == nim {
			return a, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockAlumniRepo) CreateAlumni(ctx context.Context, req model.CreateAlumniRequest) (*model.Alumni, error) {
	if m.preErr != nil {
		return nil, m.preErr
	}
	id := primitive.NewObjectID()
	now := time.Now()
	al := &model.Alumni{
		ID:         id,
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	m.alumniStore[id.Hex()] = al
	return al, nil
}

func (m *mockAlumniRepo) GetAlumniByID(ctx context.Context, id string) (*model.Alumni, error) {
	if m.preErr != nil {
		return nil, m.preErr
	}
	if a, ok := m.alumniStore[id]; ok {
		return a, nil
	}
	return nil, errors.New("not found")
}

// -------------------- "testable" service wrapper --------------------

// TestableAlumniService adalah versi service yang menerima repo via DI.
// Anda bisa membuat tipe serupa di kode production (atau refactor service agar menerima repo).
type TestableAlumniService struct {
	repo AlumniRepoMock
}

func NewTestableAlumniService(repo AlumniRepoMock) *TestableAlumniService {
	return &TestableAlumniService{repo: repo}
}

// IsAlumni implements logic yang mirip CheckAlumniService: cek nim exists
func (s *TestableAlumniService) IsAlumni(ctx context.Context, nim string) (bool, *model.Alumni, error) {
	if nim == "" {
		return false, nil, errors.New("nim empty")
	}
	al, err := s.repo.CheckAlumniByNim(ctx, nim)
	if err != nil {
		// treat not found as false without fatal error
		return false, nil, err
	}
	return true, al, nil
}

func (s *TestableAlumniService) RegisterAlumni(ctx context.Context, req model.CreateAlumniRequest) (*model.Alumni, error) {
	// minimal validation we had in original service
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return nil, errors.New("NIM, nama, jurusan, dan email wajib diisi")
	}
	return s.repo.CreateAlumni(ctx, req)
}

// -------------------- TESTS --------------------

func TestIsAlumni(t *testing.T) {
	ctx := context.Background()
	mockRepo := newMockAlumniRepo()

	// pre-populate satu alumni
	al, _ := mockRepo.CreateAlumni(ctx, model.CreateAlumniRequest{
		NIM:        "18001",
		Nama:       "Budi",
		Jurusan:    "TI",
		Angkatan:   2018,
		TahunLulus: 2022,
		Email:      "budi@example.com",
	})

	svc := NewTestableAlumniService(mockRepo)

	tests := []struct {
		name      string
		nim       string
		wantFound bool
		wantErr   bool
	}{
		{"existing nim", "18001", true, false},
		{"not exist nim", "99999", false, true}, // repo returns not found error
		{"empty nim", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, got, err := svc.IsAlumni(ctx, tt.nim)
			if (err != nil) != tt.wantErr {
				t.Fatalf("IsAlumni error = %v, wantErr %v", err, tt.wantErr)
			}
			if found != tt.wantFound {
				t.Fatalf("found = %v, want %v", found, tt.wantFound)
			}
			if found && got.NIM != al.NIM {
				t.Fatalf("expected alumni NIM %v, got %v", al.NIM, got.NIM)
			}
		})
	}
}

func TestRegisterAlumni(t *testing.T) {
	ctx := context.Background()
	mockRepo := newMockAlumniRepo()
	svc := NewTestableAlumniService(mockRepo)

	tests := []struct {
		name    string
		req     model.CreateAlumniRequest
		wantErr bool
	}{
		{
			"valid request",
			model.CreateAlumniRequest{
				NIM:        "20001",
				Nama:       "Siti",
				Jurusan:    "SI",
				Angkatan:   2020,
				TahunLulus: 2024,
				Email:      "siti@example.com",
				Password:   "secret123",
			},
			false,
		},
		{
			"missing fields",
			model.CreateAlumniRequest{
				NIM:      "",
				Nama:     "",
				Jurusan:  "",
				Email:    "",
				Password: "abc123",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.RegisterAlumni(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RegisterAlumni error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.NIM != tt.req.NIM {
				t.Fatalf("expected NIM %v, got %v", tt.req.NIM, got.NIM)
			}
		})
	}
}
