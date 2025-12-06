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
	ActivityLogRepo    pgRepo.ActivityLogRepository
	TokenRepo          TokenRepository // <-- Tambahan: Interface dari auth_service.go
}

type Services struct {
	Achievement *AchievementService
	User        *UserService
	Auth        *AuthService
	RBAC        *RBACService
	Student     *StudentService
	Lecturer    *LecturerService
	Report      *ReportService // <-- Tambahan: Report Service
}

func NewServices(db *sql.DB, mongoDB *mongodriver.Database, repos *Repos) *Services {
	// TODO: kalau mau, tambahkan logic bikin default repos kalau repos == nil

	achSvc := NewAchievementService(
		repos.AchievementRepo,
		repos.AchievementRefRepo,
		repos.StudentRepo,
		repos.UserRepo,
		repos.ActivityLogRepo,
	)

	userSvc := NewUserService(repos.UserRepo)
	
	// Update: Inject TokenRepo ke AuthService
	authSvc := NewAuthService(repos.UserRepo, repos.TokenRepo)
	
	rbacSvc := NewRBACService(repos.RolePermissionRepo, repos.PermissionRepo, repos.RoleRepo)
	studentSvc := NewStudentService(repos.StudentRepo)
	lecturerSvc := NewLecturerService(repos.LecturerRepo)

	// Tambahan: Wiring ReportService
	reportSvc := NewReportService(
		repos.AchievementRefRepo,
		repos.StudentRepo,
		repos.LecturerRepo,
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