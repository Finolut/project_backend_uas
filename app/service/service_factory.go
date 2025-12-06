package service

import (
	"database/sql"

	mongodriver "go.mongodb.org/mongo-driver/mongo"

	mongoRepo "clean-arch-copy/app/repository/mongo"
	pgRepo "clean-arch-copy/app/repository/postgre"
)

// Repos set of repo interfaces needed to create services
type Repos struct {
	UserRepo           pgRepo.UserRepository
	RoleRepo           pgRepo.RoleRepository
	PermissionRepo     pgRepo.PermissionRepository
	RolePermissionRepo pgRepo.RolePermissionRepository
	StudentRepo        pgRepo.StudentRepository
	LecturerRepo       pgRepo.LecturerRepository
	AchievementRefRepo pgRepo.AchievementRefRepository
	AchievementRepo    mongoRepo.AchievementRepository
	ActivityLogRepo    pgRepo.ActivityLogRepository // Pastikan ini ada
	TokenRepo          TokenRepository
}

type Services struct {
	Achievement *AchievementService
	User        *UserService
	Auth        *AuthService
	RBAC        *RBACService
	Student     *StudentService
	Lecturer    *LecturerService
	Report      *ReportService
}

func NewServices(db *sql.DB, mongoDB *mongodriver.Database, repos *Repos) *Services {
	// ... (kode lain tetap sama)

	achSvc := NewAchievementService(
		repos.AchievementRepo,
		repos.AchievementRefRepo,
		repos.StudentRepo,
		repos.UserRepo,
		repos.ActivityLogRepo,
	)

	userSvc := NewUserService(repos.UserRepo)
	authSvc := NewAuthService(repos.UserRepo, repos.TokenRepo)
	rbacSvc := NewRBACService(repos.RolePermissionRepo, repos.PermissionRepo, repos.RoleRepo)
	studentSvc := NewStudentService(repos.StudentRepo)
	lecturerSvc := NewLecturerService(repos.LecturerRepo)

	// Update Wiring ReportService disini:
	reportSvc := NewReportService(
		repos.AchievementRefRepo,
		repos.StudentRepo,
		repos.LecturerRepo,
		repos.ActivityLogRepo, // <-- Masukkan dependency ActivityLogRepo
	)

	return &Services{
		Achievement: achSvc,
		User:        userSvc,
		Auth:        authSvc,
		RBAC:        rbacSvc,
		Student:     studentSvc,
		Lecturer:    lecturerSvc,
		Report:      reportSvc,
	}
}
