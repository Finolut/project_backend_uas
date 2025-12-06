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
	ActivityLogRepo    pgRepo.ActivityLogRepository // <-- WAJIB ada
}

type Services struct {
	Achievement *AchievementService
	User        *UserService
	Auth        *AuthService
	RBAC        *RBACService
	Student     *StudentService
	Lecturer    *LecturerService
}

func NewServices(db *sql.DB, mongoDB *mongodriver.Database, repos *Repos) *Services {
	// TODO: kalau mau, tambahkan logic bikin default repos kalau repos == nil

	achSvc := NewAchievementService(
		repos.AchievementRepo,
		repos.AchievementRefRepo,
		repos.StudentRepo,
		repos.UserRepo,
		repos.ActivityLogRepo, // <-- argumen ke-5, ini yang diminta compiler
	)

	userSvc := NewUserService(repos.UserRepo)
	authSvc := NewAuthService(repos.UserRepo)
	rbacSvc := NewRBACService(repos.RolePermissionRepo, repos.PermissionRepo, repos.RoleRepo)
	studentSvc := NewStudentService(repos.StudentRepo)
	lecturerSvc := NewLecturerService(repos.LecturerRepo)

	return &Services{
		Achievement: achSvc,
		User:        userSvc,
		Auth:        authSvc,
		RBAC:        rbacSvc,
		Student:     studentSvc,
		Lecturer:    lecturerSvc,
	}
}
